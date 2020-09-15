package main

import (
	"context"
	"errors"
	"sync"

	grpc_peer "google.golang.org/grpc/peer"
)

type Tracker struct {
	state map[string]map[*Peer][]*Piece
	mut   sync.Mutex
}

func NewTracker() *Tracker {
	return &Tracker{
		state: make(map[string]map[*Peer][]*Piece),
	}
}

func (tracker *Tracker) Announce(ctx context.Context, req *AnnounceRequest) (*AnnounceReply, error) {
	p, ok := grpc_peer.FromContext(ctx)
	if !ok {
		return nil, errors.New(RegisterPeerContextNotOkError)
	}

	peer := &Peer{Address: p.Addr.String()}

	tracker.mut.Lock()
	_, ok = tracker.state[req.TorrentSha256]
	if !ok {
		tracker.state[req.TorrentSha256] = make(map[*Peer][]*Piece)
	}
	tracker.state[req.TorrentSha256][peer] = req.TorrentPieces
	peers := make([]*Peer, 0)
	for p := range tracker.state[req.TorrentSha256] {
		if p.Address == peer.Address {
			continue
		}

		peers = append(peers, p)
	}
	tracker.mut.Unlock()

	return &AnnounceReply{
		Peers: peers,
	}, nil
}
