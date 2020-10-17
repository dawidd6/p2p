package daemon

import (
	"context"
	"net"
	"path/filepath"
	"sync"
	"time"

	"github.com/dawidd6/p2p/pkg/errors"

	"github.com/dawidd6/p2p/pkg/utils"

	"github.com/dawidd6/p2p/pkg/torrent"

	"github.com/dawidd6/p2p/pkg/tracker"

	"google.golang.org/grpc"
)

type Config struct {
	DownloadsDir string
	Address      string
	SeedAddress  string
}

type Daemon struct {
	// file_sha256 piece_number torrent
	torrents map[string]
	config   *Config
	mut      sync.Mutex
	UnimplementedDaemonServer
	UnimplementedSeedServer
}

func New(config *Config) *Daemon {
	return &Daemon{
		torrents: make(map[string]map[uint64]*torrent.Torrent),
		config:   config,
	}
}

func (daemon *Daemon) Listen() error {
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
		RegisterSeedServer(server, daemon)
		ch <- server.Serve(listener)
	}()

	return <-ch
}

func (daemon *Daemon) getPeerAddresses(t *torrent.Torrent) ([]string, error) {
	peers := make([]string, 0)

	for _, url := range t.TrackerUrls {
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

	for i, pieceHash := range t.PieceHashes {
		// TODO, throttle for now
		<-time.After(time.Second)

		request := &FetchRequest{
			FileHash:    t.FileHash,
			PieceNumber: uint64(i),
		}

		for _, peerAddr := range peerAddresses {
			// TODO use DialContext everywhere
			conn, err := grpc.DialContext(context.TODO(), peerAddr, grpc.WithInsecure())
			if err != nil {
				return err
			}

			response, err := NewSeedClient(conn).Fetch(context.TODO(), request)
			if err != nil {
				return err
			}

			// Peer does not have this piece, ask someone else
			if response.PieceData == nil {
				continue
			}

			// Got the piece wrongly, ask someone else
			if utils.Sha256Sum(response.PieceData) != pieceHash {
				continue
			}

			err = utils.WriteFilePiece(filepath.Join(daemon.config.DownloadsDir, t.FileName), uint64(i), response.PieceData)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (daemon *Daemon) add(t *torrent.Torrent) {

}

func (daemon *Daemon) filePath(t *torrent.Torrent) string {
	return filepath.Join(daemon.config.DownloadsDir, t.FileName)
}

func (daemon *Daemon) Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error) {
	daemon.mut.Lock()
	t, ok := daemon.torrents[req.FileHash]
	daemon.mut.Unlock()
	if !ok {
		return nil, errors.TorrentNotFound
	}

	piece, err := utils.ReadFilePiece(daemon.filePath(t), t.PieceSize, req.PieceNumber)
	if err != nil {
		return nil, err
	}

	return &FetchResponse{PieceData: piece}, nil
}

func (daemon *Daemon) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	if utils.FileExists(daemon.filePath(req.Torrent)) {

	}

	return &AddResponse{Torrent: req.Torrent}, nil
}
