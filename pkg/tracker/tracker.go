package tracker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dawidd6/p2p/pkg/proto"
)

type Tracker struct {
	// torrent_sha256 peer_address peer_announce_timestamp
	index map[string]map[string]uint64
	mut   sync.Mutex
	proto.UnimplementedTrackerServer
}

func NewTracker() *Tracker {
	return &Tracker{
		index: make(map[string]map[string]uint64),
	}
}

func (tracker *Tracker) GoClean() {
	go func() {
		for {
			<-time.After(time.Second * 30)
			log.Println("Clean", "start")

			tracker.mut.Lock()
			currentTimestamp := uint64(time.Now().UTC().Unix())

			for torrentSha256 := range tracker.index {
				for peerAddress, peerTimestamp := range tracker.index[torrentSha256] {
					if currentTimestamp-peerTimestamp > 60 {
						log.Println("Clean", torrentSha256, peerAddress, currentTimestamp, peerTimestamp)
						delete(tracker.index[torrentSha256], peerAddress)
					}
				}
				if len(tracker.index[torrentSha256]) == 0 {
					log.Println("Clean", torrentSha256, len(tracker.index[torrentSha256]))
					delete(tracker.index, torrentSha256)
				}
			}
			tracker.mut.Unlock()

			log.Println("Clean", "stop")
		}
	}()
}

func (tracker *Tracker) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterReply, error) {
	log.Println("Announce", req.PeerAddress, req.TorrentSha256)

	tracker.mut.Lock()
	if !tracker.isTorrentSha256Present(req.TorrentSha256) {
		tracker.index[req.TorrentSha256] = make(map[string]uint64)
	}
	tracker.index[req.TorrentSha256][req.PeerAddress] = uint64(time.Now().UTC().Unix())
	peerAddresses := tracker.peerAddressesForTorrentSha256(req.TorrentSha256, req.PeerAddress)
	tracker.mut.Unlock()

	return &proto.RegisterReply{
		PeerAddresses: peerAddresses,
	}, nil
}

func (tracker *Tracker) List(ctx context.Context, req *proto.ListRequest) (*proto.ListReply, error) {
	tracker.mut.Lock()
	peerAddresses := tracker.peerAddressesForTorrentSha256(req.TorrentSha256, req.PeerAddress)
	tracker.mut.Unlock()

	return &proto.ListReply{
		PeerAddresses: peerAddresses,
	}, nil
}

func (tracker *Tracker) isTorrentSha256Present(torrentSha256 string) bool {
	_, ok := tracker.index[torrentSha256]

	return ok
}

func (tracker *Tracker) peerAddressesForTorrentSha256(torrentSha256, excludedAddress string) []string {
	peerAddresses := make([]string, 0)

	if !tracker.isTorrentSha256Present(torrentSha256) {
		return peerAddresses
	}

	for address := range tracker.index[torrentSha256] {
		if address == excludedAddress {
			continue
		}

		peerAddresses = append(peerAddresses, address)
	}

	return peerAddresses
}
