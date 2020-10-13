package daemon

import (
	"context"
	"log"
	"net"

	"github.com/dawidd6/p2p/pkg/seed"

	"github.com/dawidd6/p2p/pkg/torrent"
	"github.com/dawidd6/p2p/pkg/tracker"

	"google.golang.org/grpc"
)

const DefaultListenAddr = "localhost:8888"

type Daemon struct {
	torrents []*torrent.Torrent
	address  string
	UnimplementedDaemonServer
}

func New(address string) *Daemon {
	return &Daemon{
		torrents: make([]*torrent.Torrent, 0),
		address:  address,
	}
}

func (daemon *Daemon) Listen() error {
	listener, err := net.Listen("tcp", daemon.address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterDaemonServer(server, daemon)
	return server.Serve(listener)
}

func (daemon *Daemon) Add(ctx context.Context, in *AddRequest) (*AddResponse, error) {
	peers := &tracker.RegisterReply{}

	for _, url := range in.Torrent.Trackers {
		conn, err := grpc.Dial(url, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		request := &tracker.RegisterRequest{
			TorrentSha256: in.Torrent.Sha256,
			PeerAddress:   daemon.address,
		}
		client := tracker.NewTrackerClient(conn)
		peers, err = client.Register(context.TODO(), request)
		if err != nil {
			return nil, err
		}
	}

	log.Println("Add")

	request := &seed.GetRequest{
		TorrentFileSha256: in.Torrent.Sha256,
		PieceNumber:       1, // TODO
	}

	for _, peerAddr := range peers.PeerAddresses {
		conn, err := grpc.Dial(peerAddr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		client := seed.NewSeedClient(conn)
		_, err = client.Get(context.TODO(), request)
		if err != nil {
			return nil, err
		}
	}

	return &AddResponse{Torrent: in.Torrent}, nil
}
