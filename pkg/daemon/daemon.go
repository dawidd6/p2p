package daemon

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/dawidd6/p2p/pkg/hash"
	"github.com/dawidd6/p2p/pkg/piece"

	"github.com/dawidd6/p2p/pkg/file"

	"github.com/dawidd6/p2p/pkg/mmap"

	"github.com/dawidd6/p2p/pkg/errors"
	"github.com/dawidd6/p2p/pkg/torrent"
	"github.com/dawidd6/p2p/pkg/tracker"
	"google.golang.org/grpc"
)

type Config struct {
	DownloadsDir string
	Address      string
	SeedAddress  string
}

type State struct {
	DownloadedPieceNumbers *mmap.MMap
	File                   *os.File
	Torrent                *torrent.Torrent
	AnnounceChannel        chan struct{}
}

type Daemon struct {
	torrents *mmap.MMap // map[string(fileHash)]*State
	config   *Config
	UnimplementedDaemonServer
	UnimplementedSeederServer
}

func Run(config *Config) error {
	daemon := &Daemon{
		torrents: mmap.New(),
		config:   config,
	}

	ch := make(chan error)

	go func() {
		listener, err := net.Listen("tcp", daemon.config.Address)
		if err != nil {
			ch <- err
			return
		}

		server := grpc.NewServer()
		RegisterDaemonServer(server, daemon)
		ch <- server.Serve(listener)
	}()

	go func() {
		listener, err := net.Listen("tcp", daemon.config.SeedAddress)
		if err != nil {
			ch <- err
			return
		}

		server := grpc.NewServer()
		RegisterSeederServer(server, daemon)
		ch <- server.Serve(listener)
	}()

	return <-ch
}

func (daemon *Daemon) announce(state *State) (map[string]*tracker.AnnounceResponse, error) {
	track := make(map[string]*tracker.AnnounceResponse)

	for _, trackerAddr := range state.Torrent.TrackerAddresses {
		conn, err := grpc.Dial(trackerAddr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		request := &tracker.AnnounceRequest{
			FileHash:    state.Torrent.FileHash,
			PeerAddress: daemon.config.SeedAddress,
		}

		response, err := tracker.NewTrackerClient(conn).Announce(context.TODO(), request)
		if err != nil {
			return nil, err
		}
		track[trackerAddr] = response

	}

	return track, nil
}

func (daemon *Daemon) fetch(state *State) {
	err := file.CreateDirectory(daemon.config.DownloadsDir)

	if err != nil {
		log.Println(err)
		return
	}

	a, err := daemon.announce(state)
	if err != nil {
		log.Println(err)
	}

	go func() {
		for {
			select {
			case <-time.After(time.Duration(a) * time.Second):
				peerAddresses, announceInterval, err = daemon.announce(state)
			case <-state.AnnounceChannel:
				break
			}
		}
	}()

	state.File, err = file.Open(daemon.config.DownloadsDir, state.Torrent.FileName)
	if err != nil {
		log.Println(err)
		return
	}

	err = torrent.Verify(state.Torrent, daemon.config.DownloadsDir)
	if err == nil {
		log.Println("seeding")
		return
	}

	for {
		if len(peerAddresses) == 0 {
			log.Println("no peers available")
			<-time.After(time.Second * 5)
			continue
		}

		for i, pieceHash := range state.Torrent.PieceHashes {
			pieceNumber := uint64(i)

			// TODO, throttle for now
			<-time.After(time.Second * 1)

			request := &SeedRequest{
				FileHash:    state.Torrent.FileHash,
				PieceNumber: pieceNumber,
			}

			pieceData, err := piece.Read(state.File, state.Torrent.PieceSize, pieceNumber)
			if err != nil {
				log.Println(err)
			}

			if hash.Compute(pieceData) == pieceHash {
				log.Println(pieceHash, "already downloaded piece")
				continue
			}

			for _, peerAddr := range peerAddresses {
				log.Println("Fetch", pieceNumber, pieceHash)

				// TODO use DialContext everywhere
				conn, err := grpc.DialContext(context.TODO(), peerAddr, grpc.WithInsecure())
				if err != nil {
					log.Println(err)
					continue
				}

				response, err := NewSeederClient(conn).Seed(context.TODO(), request)
				if err != nil {
					log.Println(err)
					continue
				}

				// Peer does not have this piece, ask someone else
				if response.PieceData == nil {
					log.Println(peerAddr, "no piece")
					continue
				}

				// Got the piece wrongly, ask someone else
				if hash.Compute(response.PieceData) != pieceHash {
					log.Println(peerAddr, "wrong piece")
					continue
				}

				err = piece.Write(state.File, pieceNumber, pieceData)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}

		err = torrent.Verify(state.Torrent, daemon.config.DownloadsDir)
		if err != nil {
			log.Println(err)
			continue
		} else {
			log.Println("completed")
			return
		}
	}
}

func (daemon *Daemon) filePath(t *torrent.Torrent) string {
	return filepath.Join(daemon.config.DownloadsDir, t.FileName)
}

func (daemon *Daemon) Seed(ctx context.Context, req *SeedRequest) (*SeedResponse, error) {
	if !daemon.torrents.Has(req.FileHash) {
		return nil, errors.TorrentNotFound
	}

	state := daemon.torrents.Get(req.FileHash).(State)

	pieceData, err := piece.Read(state.File, state.Torrent.PieceSize, req.PieceNumber)
	if err != nil {
		return nil, err
	}

	return &SeedResponse{PieceData: pieceData}, nil
}

func (daemon *Daemon) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	if daemon.torrents.Has(req.Torrent.FileHash) {
		return nil, errors.TorrentAlreadyAdded
	}

	state := &State{
		DownloadedPieceNumbers: mmap.New(),
		Torrent:                req.Torrent,
		AnnounceChannel:        make(chan struct{}),
	}

	daemon.torrents.Set(req.Torrent.FileHash, state)

	go daemon.fetch(state)

	return &AddResponse{Torrent: req.Torrent}, nil
}
