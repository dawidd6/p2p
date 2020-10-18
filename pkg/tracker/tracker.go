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
	Address          string
	AnnounceInterval time.Duration
	CleanInterval    time.Duration
}

type Tracker struct {
	// torrent_sha256 peer_address peer_announce_timestamp
	index  map[string]map[string]uint64
	mut    sync.Mutex
	config *Config
	UnimplementedTrackerServer
}

func Run(config *Config) error {
	tracker := &Tracker{
		index:  make(map[string]map[string]uint64),
		config: config,
	}

	go tracker.clean()

	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterTrackerServer(server, tracker)
	return server.Serve(listener)
}

func (tracker *Tracker) clean() {
	for {
		<-time.After(tracker.config.CleanInterval)

		currentTimestamp := uint64(time.Now().UTC().Unix())

		tracker.mut.Lock()
		for fileHash := range tracker.index {
			for peerAddress, peerTimestamp := range tracker.index[fileHash] {
				if currentTimestamp-peerTimestamp > 60 {
					log.Println("Clean", fileHash, peerAddress)
					delete(tracker.index[fileHash], peerAddress)
				}
			}
			if len(tracker.index[fileHash]) == 0 {
				log.Println("Clean", fileHash)
				delete(tracker.index, fileHash)
			}
		}
		tracker.mut.Unlock()
	}
}

func (tracker *Tracker) Announce(ctx context.Context, req *AnnounceRequest) (*AnnounceResponse, error) {
	log.Println("Announce", req.PeerAddress, req.FileHash)
	peerAddresses := make([]string, 0)

	tracker.mut.Lock()
	if _, ok := tracker.index[req.FileHash]; !ok {
		tracker.index[req.FileHash] = make(map[string]uint64)
	}
	tracker.index[req.FileHash][req.PeerAddress] = uint64(time.Now().UTC().Unix())
	for address := range tracker.index[req.FileHash] {
		if address == req.PeerAddress {
			continue
		}

		peerAddresses = append(peerAddresses, address)
	}
	tracker.mut.Unlock()

	return &AnnounceResponse{
		PeerAddresses:    peerAddresses,
		AnnounceInterval: uint32(tracker.config.AnnounceInterval.Seconds()),
	}, nil
}
