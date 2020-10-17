package daemon

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"sort"
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

type Daemon struct {
	torrents []*torrent.Torrent
	mut      sync.Mutex
	config   *Config
	UnimplementedDaemonServer
	UnimplementedSeederServer
}

func Run(config *Config) error {
	daemon := &Daemon{
		torrents: make([]*torrent.Torrent, 0),
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

func (daemon *Daemon) fetch(t *torrent.Torrent) error {
	peerAddresses, err := daemon.getPeerAddresses(t)
	if err != nil {
		return err
	}

	err = os.MkdirAll(daemon.config.DownloadsDir, os.ModeDir)
	if err != nil {
		return err
	}

	err = utils.AllocateZeroedFile(daemon.filePath(t), t.FileSize)
	if err != nil {
		return err
	}

	for i, pieceHash := range t.PieceHashes {
		// TODO, throttle for now
		<-time.After(time.Second)

		request := &SeedRequest{
			FileHash:    t.FileHash,
			PieceNumber: uint64(i),
		}

		for _, peerAddr := range peerAddresses {
			// TODO use DialContext everywhere
			conn, err := grpc.DialContext(context.TODO(), peerAddr, grpc.WithInsecure())
			if err != nil {
				continue
			}

			response, err := NewSeederClient(conn).Seed(context.TODO(), request)
			if err != nil {
				continue
			}

			// Peer does not have this piece, ask someone else
			if response.PieceData == nil {
				continue
			}

			// Got the piece wrongly, ask someone else
			if utils.Sha256Sum(response.PieceData) != pieceHash {
				continue
			}

			err = utils.WriteFilePiece(daemon.filePath(t), uint64(i), response.PieceData)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

func (daemon *Daemon) filePath(t *torrent.Torrent) string {
	return filepath.Join(daemon.config.DownloadsDir, t.FileName)
}

func (daemon *Daemon) Seed(ctx context.Context, req *SeedRequest) (*SeedResponse, error) {
	daemon.mut.Lock()
	n := sort.Search(len(daemon.torrents), func(i int) bool {
		if daemon.torrents[i].FileHash == req.FileHash {
			return true
		}
		return false
	})
	if n < len(daemon.torrents) {
	} else {
		return nil, errors.TorrentNotFound
	}
	t := daemon.torrents[n]
	daemon.mut.Unlock()

	piece, err := utils.ReadFilePiece(daemon.filePath(t), t.PieceSize, req.PieceNumber)
	if err != nil {
		return nil, err
	}

	return &SeedResponse{PieceData: piece}, nil
}

func (daemon *Daemon) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	daemon.mut.Lock()
	n := sort.Search(len(daemon.torrents), func(i int) bool {
		if daemon.torrents[i].FileHash == req.Torrent.FileHash {
			return true
		}
		return false
	})
	if n < len(daemon.torrents) {
		return nil, errors.TorrentAlreadyAdded
	} else {
		daemon.torrents = append(daemon.torrents, req.Torrent)
	}
	daemon.mut.Unlock()

	go func() {
		err := daemon.fetch(req.Torrent)
		if err != nil {
			log.Println("Fetch", err)
		}
	}()

	return &AddResponse{Torrent: req.Torrent}, nil
}
