package tracker

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/dawidd6/p2p/pkg/defaults"

	"google.golang.org/grpc"
)

type Config struct {
	ListenAddress    string
	AnnounceInterval time.Duration
	CleanInterval    time.Duration
}

type Tracker struct {
	config *Config
	index  map[string]map[string]time.Time
	mutex  sync.Mutex
	UnimplementedTrackerServer
}

func Run(config *Config) error {
	tracker := &Tracker{
		config: config,
		index:  make(map[string]map[string]time.Time),
	}

	go func() {
		for {
			<-time.After(config.CleanInterval)
			tracker.clean()
		}
	}()

	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		return err
	}

	server := grpc.NewServer(grpc.ConnectionTimeout(defaults.TrackerConnTimeout))
	RegisterTrackerServer(server, tracker)
	return server.Serve(listener)
}

func (tracker *Tracker) clean() {
	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()

	for fileHash := range tracker.index {
		for peerAddress, peerTimestamp := range tracker.index[fileHash] {
			if time.Since(peerTimestamp) > tracker.config.AnnounceInterval*2 {
				log.Println("Clean", fileHash, peerAddress)
				delete(tracker.index[fileHash], peerAddress)
			}
			if len(tracker.index[fileHash]) == 0 {
				log.Println("Clean", fileHash)
				delete(tracker.index, fileHash)
			}
		}
	}
}

func (tracker *Tracker) Announce(ctx context.Context, req *AnnounceRequest) (*AnnounceResponse, error) {
	log.Println("Announce", req.PeerAddress)

	tracker.mutex.Lock()
	defer tracker.mutex.Unlock()

	if _, ok := tracker.index[req.FileHash]; !ok {
		tracker.index[req.FileHash] = make(map[string]time.Time)
	}

	i := 0
	peerAddresses := make([]string, len(tracker.index[req.FileHash]))
	for peerAddress := range tracker.index[req.FileHash] {
		peerAddresses[i] = peerAddress
		i++
	}

	tracker.index[req.FileHash][req.PeerAddress] = time.Now()

	return &AnnounceResponse{
		PeerAddresses:    peerAddresses,
		AnnounceInterval: uint32(tracker.config.AnnounceInterval.Seconds()),
	}, nil
}
