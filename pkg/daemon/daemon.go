package daemon

import (
	"context"
	"log"

	"github.com/dawidd6/p2p/pkg/proto"

	"google.golang.org/grpc"
)

type Daemon struct {
	torrents []*proto.Torrent
	address  string
	proto.UnimplementedDaemonServer
}

func NewDaemon(address string) *Daemon {
	return &Daemon{
		torrents: make([]*proto.Torrent, 0),
		address:  address,
	}
}

func (daemon *Daemon) Add(ctx context.Context, in *proto.AddRequest) (*proto.AddResponse, error) {
	for _, url := range in.Torrent.Trackers {
		conn, err := grpc.Dial(url, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		request := &proto.RegisterRequest{
			TorrentSha256: in.Torrent.Sha256,
			PeerAddress:   daemon.address,
		}
		client := proto.NewTrackerClient(conn)
		_, err = client.Register(context.TODO(), request)
		if err != nil {
			return nil, err
		}
	}

	log.Println("Add")

	return &proto.AddResponse{Torrent: in.Torrent}, nil
}
