package daemon

import (
	"context"

	"github.com/dawidd6/p2p/proto"
)

type Daemon struct {
	torrents []*proto.Torrent
}

func NewDaemon() *Daemon {
	return &Daemon{
		torrents: make([]*proto.Torrent, 0),
	}
}

func (daemon *Daemon) Add(ctx context.Context, in *proto.AddRequest) (*proto.AddResponse, error) {
	return &proto.AddResponse{Torrent: in.Torrent}, nil
}

func (daemon *Daemon) Delete(ctx context.Context, in *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	return &proto.DeleteResponse{Torrent: in.Torrent}, nil
}
