package main

import (
	"context"
	"errors"
	"sync"

	"google.golang.org/grpc/peer"
)

type Tracker struct {
	state map[string]map[*Peer][]Piece
	mut   sync.Mutex
}

func NewTracker() *Tracker {
	return &Tracker{}
}

func (tracker *Tracker) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	pp, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New(RegisterPeerContextNotOkError)
	}
	p := &Peer{Address: pp.Addr.String()}
	tracker.mut.Lock()
	_, ok = tracker.state[in.TorrentSha256][p]
	tracker.mut.Unlock()
	return &RegisterResponse{}, nil
}
