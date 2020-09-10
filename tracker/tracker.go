package tracker

import (
	"context"
	"sync"

	"github.com/dawidd6/p2p/proto"

	grpeer "google.golang.org/grpc/peer"
)

type Tracker struct {
	peersToMetadata map[*proto.Peer][]*proto.Torrent
	mut             sync.Mutex
}

func NewTracker() *Tracker {
	return &Tracker{
		peersToMetadata: make(map[*proto.Peer][]*proto.Torrent, 0),
	}
}

func (tracker *Tracker) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	pp, _ := grpeer.FromContext(ctx)
	p := &proto.Peer{Address: pp.Addr.String()}
	tracker.mut.Lock()
	tracker.peersToMetadata[p] = in.Torrents
	tracker.mut.Unlock()
	return &proto.RegisterResponse{}, nil
}
