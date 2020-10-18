package daemon

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

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
	DownloadedPieceNumbers []uint64
}

type Daemon struct {
	torrents map[*torrent.Torrent]*State
	mut      sync.Mutex
	config   *Config
	UnimplementedDaemonServer
	UnimplementedSeederServer
}

func Run(config *Config) error {
	daemon := &Daemon{
		torrents: make(map[*torrent.Torrent]*State),
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

func (daemon *Daemon) getPeerAddresses(t *torrent.Torrent) ([]string, error) {
	peers := make([]string, 0)

	for _, url := range t.TrackerAddresses {
		conn, err := grpc.Dial(url, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		request := &tracker.AnnounceRequest{
			FileHash:    t.FileHash,
			PeerAddress: daemon.config.SeedAddress,
		}

		response, err := tracker.NewTrackerClient(conn).Announce(context.TODO(), request)
		if err != nil {
			return nil, err
		}

		// Add peer address if it does not already exist
		for _, peer := range response.PeerAddresses {
			exists := false

			for i := range peers {
				if peers[i] == peer {
					exists = true
				}
			}

			if !exists {
				peers = append(peers, peer)
			}
		}
	}

	return peers, nil
}

func (daemon *Daemon) fetch(t *torrent.Torrent) {
	err := os.MkdirAll(daemon.config.DownloadsDir, 0775)
	if err != nil {
		log.Println(err)
	}

	peerAddresses, err := daemon.getPeerAddresses(t)
	if err != nil {
		log.Println(err)
	}

	go func() {
		for {
			<-time.After(time.Second * 30)
			peerAddresses, err = daemon.getPeerAddresses(t)
		}
	}()

	err = torrent.Verify(t, daemon.config.DownloadsDir)
	if err == nil {
		log.Println("Fetch", "seeding")
		return
	} else {
		err = utils.AllocateZeroedFile(daemon.filePath(t), t.FileSize)
		if err != nil {
			log.Println(err)
		}
	}

	for {
		if len(peerAddresses) == 0 {
			log.Println("no peers available")
			<-time.After(time.Second * 5)
			continue
		}

		for i, pieceHash := range t.PieceHashes {
			// TODO, throttle for now
			<-time.After(time.Second * 1)

			request := &SeedRequest{
				FileHash:    t.FileHash,
				PieceNumber: uint64(i),
			}

			piece, err := utils.ReadFilePiece(daemon.filePath(t), t.PieceSize, uint64(i))
			if err != nil {
				log.Println(err)
				continue
			}

			if utils.Sha256Sum(piece) == pieceHash {
				log.Println(pieceHash, "already downloaded piece")
				continue
			}

			for _, peerAddr := range peerAddresses {
				log.Println("Fetch", i, pieceHash)

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

				err = utils.WriteFilePiece(daemon.filePath(t), uint64(i), response.PieceData)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}

		err = torrent.Verify(t, daemon.config.DownloadsDir)
		if err != nil {
			log.Println(err)
			continue
		} else {
			log.Println("complete")
			break
		}
	}
}

func (daemon *Daemon) filePath(t *torrent.Torrent) string {
	return filepath.Join(daemon.config.DownloadsDir, t.FileName)
}

func (daemon *Daemon) matchTorrent(fileHash string) *torrent.Torrent {
	daemon.mut.Lock()
	defer daemon.mut.Unlock()
	for t := range daemon.torrents {
		if t.FileHash == fileHash {
			return t
		}
	}
	return nil
}

func (daemon *Daemon) Seed(ctx context.Context, req *SeedRequest) (*SeedResponse, error) {
	t := daemon.matchTorrent(req.FileHash)
	if t == nil {
		return nil, errors.TorrentNotFound
	}

	piece, err := utils.ReadFilePiece(daemon.filePath(t), t.PieceSize, req.PieceNumber)
	if err != nil {
		return nil, err
	}

	if utils.Sha256Sum(piece) != t.PieceHashes[req.PieceNumber] {
		return nil, errors.PieceChecksumMismatchError
	}

	return &SeedResponse{PieceData: piece}, nil
}

func (daemon *Daemon) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	t := daemon.matchTorrent(req.Torrent.FileHash)
	if t != nil {
		return nil, errors.TorrentAlreadyAdded
	} else {
		daemon.mut.Lock()
		daemon.torrents[req.Torrent] = &State{}
		daemon.mut.Unlock()
	}

	go daemon.fetch(req.Torrent)

	return &AddResponse{Torrent: req.Torrent}, nil
}
