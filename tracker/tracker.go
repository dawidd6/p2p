package tracker

import (
	"context"

	"github.com/dawidd6/p2p/torrent"

	grpeer "google.golang.org/grpc/peer"

	"github.com/dawidd6/p2p/peer"
)

type Tracker struct {
	peersToMetadata map[*peer.Peer][]*torrent.Torrent
}

func NewTracker() *Tracker {
	return &Tracker{}
}

func (tracker *Tracker) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	pp, _ := grpeer.FromContext(ctx)
	p := &peer.Peer{Address: pp.Addr.String()}
	tracker.peersToMetadata[p] = in.Torrents
	return &RegisterResponse{}, nil
}

func (tracker *Tracker) Lookup(ctx context.Context, in *LookupRequest) (*LookupResponse, error) {
	return &LookupResponse{}, nil
}
