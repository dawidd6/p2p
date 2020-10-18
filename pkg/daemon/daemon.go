package daemon

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/dawidd6/p2p/pkg/mmap"

	"github.com/dawidd6/p2p/pkg/errors"
	"github.com/dawidd6/p2p/pkg/torrent"
	"github.com/dawidd6/p2p/pkg/tracker"
	"github.com/dawidd6/p2p/pkg/utils"

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

func (daemon *Daemon) announce(state *State) ([]string, uint32, error) {
	peers := mmap.New()
	announceInterval := uint32(0)

	for _, trackerAddr := range state.Torrent.TrackerAddresses {
		conn, err := grpc.Dial(trackerAddr, grpc.WithInsecure())
		if err != nil {
			return nil, 0, err
		}

		request := &tracker.AnnounceRequest{
			FileHash:    state.Torrent.FileHash,
			PeerAddress: daemon.config.SeedAddress,
		}

		response, err := tracker.NewTrackerClient(conn).Announce(context.TODO(), request)
		if err != nil {
			return nil, 0, err
		}

		for _, peer := range response.PeerAddresses {
			peers.Set(peer, struct{}{})
		}

		announceInterval = response.AnnounceInterval
	}

	return peers.Keys().([]string), announceInterval, nil
}

func (daemon *Daemon) fetch(state *State) {
	err := utils.CreateDir(daemon.config.DownloadsDir)
	if err != nil {
		log.Println(err)
		return
	}

	peerAddresses, announceInterval, err := daemon.announce(state)
	if err != nil {
		log.Println(err)
	}

	go func() {
		for {
			select {
			case <-time.After(time.Duration(announceInterval) * time.Second):
				peerAddresses, announceInterval, err = daemon.announce(state)
			case <-state.AnnounceChannel:
				break
			}
		}
	}()

	file, err := utils.OpenFile(daemon.config.DownloadsDir, state.Torrent.FileName)
	if err != nil {
		log.Println(err)
		return
	}

	state.File = file

	fileContent, err := utils.ReadFile(file)
	if err != nil {
		log.Println(err)
		return
	}

	if utils.Sha256Sum(fileContent) == state.Torrent.FileHash {
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
			pieceOffset := int64(state.Torrent.PieceSize * pieceNumber)

			// TODO, throttle for now
			<-time.After(time.Second * 1)

			request := &SeedRequest{
				FileHash:    state.Torrent.FileHash,
				PieceNumber: pieceNumber,
			}

			piece := make([]byte, state.Torrent.PieceSize)
			n, err := file.ReadAt(piece, pieceOffset)
			if err != nil && err != io.EOF {
				log.Println(err)
			}
			piece = piece[:n]

			if utils.Sha256Sum(piece) == pieceHash {
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
				if utils.Sha256Sum(response.PieceData) != pieceHash {
					log.Println(peerAddr, "wrong piece")
					continue
				}

				_, err = file.WriteAt(response.PieceData, pieceOffset)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}

		fileContent, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
		}

		if utils.Sha256Sum(fileContent) != state.Torrent.FileHash {
			log.Println(errors.FileChecksumMismatchError)
			continue
		} else {
			log.Println("completed")
			break
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

	piece := make([]byte, state.Torrent.PieceSize)
	n, err := state.File.ReadAt(piece, int64(state.Torrent.PieceSize*req.PieceNumber))
	if err != nil && err != io.EOF {
		log.Println(err)
	}
	piece = piece[:n]

	if utils.Sha256Sum(piece) != state.Torrent.PieceHashes[req.PieceNumber] {
		return nil, errors.PieceChecksumMismatchError
	}

	return &SeedResponse{PieceData: piece}, nil
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
