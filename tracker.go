package main

import (
	"context"
	"errors"
	"log"
	"sync"

	"google.golang.org/grpc/metadata"
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New(MetadataContextNotOkError)
	}

	address := md.Get(":authority")
	peer := &Peer{Address: address[0]}

	log.Println("Announce", peer.Address)

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
