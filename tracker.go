package main

import (
	"context"
	"sync"

	"google.golang.org/grpc/peer"
)

type Tracker struct {
	peersToMetadata map[*Peer][]*Torrent
	mut             sync.Mutex
}

func NewTracker() *Tracker {
	return &Tracker{
		peersToMetadata: make(map[*Peer][]*Torrent, 0),
	}
}

func (tracker *Tracker) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	pp, _ := peer.FromContext(ctx)
	p := &Peer{Address: pp.Addr.String()}
	tracker.mut.Lock()
	tracker.peersToMetadata[p] = in.Torrents
	tracker.mut.Unlock()
	return &RegisterResponse{}, nil
}
