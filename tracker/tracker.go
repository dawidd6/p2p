package tracker

import (
	"context"
	"log"
	"sync"

	"github.com/dawidd6/p2p/errors"
	"github.com/dawidd6/p2p/peer"
	piece "github.com/dawidd6/p2p/piece"

	"google.golang.org/grpc/metadata"
)

type Tracker struct {
	state map[string]map[*peer.Peer][]*piece.Piece
	mut   sync.Mutex
}

func NewTracker() *Tracker {
	return &Tracker{
		state: make(map[string]map[*peer.Peer][]*piece.Piece),
	}
}

func (tracker *Tracker) Announce(ctx context.Context, req *AnnounceRequest) (*AnnounceReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.MetadataContextNotOkError
	}

	address := md.Get(":authority")
	p := &peer.Peer{Address: address[0]}

	log.Println("Announce", p.Address)

	tracker.mut.Lock()
	_, ok = tracker.state[req.TorrentSha256]
	if !ok {
		tracker.state[req.TorrentSha256] = make(map[*peer.Peer][]*piece.Piece)
	}
	tracker.state[req.TorrentSha256][p] = req.TorrentPieces
	peers := make([]*peer.Peer, 0)
	for p := range tracker.state[req.TorrentSha256] {
		if p.Address == p.Address {
			continue
		}

		peers = append(peers, p)
	}
	tracker.mut.Unlock()

	return &AnnounceReply{
		Peers: peers,
	}, nil
}
