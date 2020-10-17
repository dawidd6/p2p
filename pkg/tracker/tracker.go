package tracker

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type Config struct {
	Address string
}

type Tracker struct {
	// torrent_sha256 peer_address peer_announce_timestamp
	index  map[string]map[string]uint64
	mut    sync.Mutex
	config *Config
	UnimplementedTrackerServer
}

func New(config *Config) *Tracker {
	return &Tracker{
		index:  make(map[string]map[string]uint64),
		config: config,
	}
}

func (tracker *Tracker) Listen() error {
	listener, err := net.Listen("tcp", tracker.config.Address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterTrackerServer(server, tracker)
	return server.Serve(listener)
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

func (tracker *Tracker) Announce(ctx context.Context, req *AnnounceRequest) (*AnnounceResponse, error) {
	log.Println("Announce", req.PeerAddress, req.FileHash)

	tracker.mut.Lock()
	if !tracker.isHashPresent(req.FileHash) {
		tracker.index[req.FileHash] = make(map[string]uint64)
	}
	tracker.index[req.FileHash][req.PeerAddress] = uint64(time.Now().UTC().Unix())
	peerAddresses := tracker.peerAddressesForHash(req.FileHash, req.PeerAddress)
	tracker.mut.Unlock()

	return &AnnounceResponse{
		PeerAddresses: peerAddresses,
	}, nil
}

func (tracker *Tracker) isHashPresent(torrentSha256 string) bool {
	_, ok := tracker.index[torrentSha256]

	return ok
}

func (tracker *Tracker) peerAddressesForHash(torrentSha256, excludedAddress string) []string {
	peerAddresses := make([]string, 0)

	if !tracker.isHashPresent(torrentSha256) {
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
