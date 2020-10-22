package tracker

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

const (
	ListenAddress    = "127.0.0.1:8889"
	AnnounceInterval = time.Second * 30
	CleanInterval    = time.Second * 60
	ConnTimeout      = time.Second * 5
)

type Config struct {
	ListenAddress    string
	AnnounceInterval time.Duration
	CleanInterval    time.Duration
}

type Tracker struct {
	config  *Config
	index   map[string]map[string]time.Time // fileHash peerAddress peerTimestamp
	mutex   sync.RWMutex
	cleaner *time.Ticker
	UnimplementedTrackerServer
}

func Run(config *Config) error {
	tracker := &Tracker{
		config:  config,
		index:   make(map[string]map[string]time.Time),
		cleaner: time.NewTicker(config.CleanInterval),
	}

	go tracker.clean()

	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		return err
	}

	server := grpc.NewServer(grpc.ConnectionTimeout(ConnTimeout))
	RegisterTrackerServer(server, tracker)
	return server.Serve(listener)
}

func (tracker *Tracker) clean() {
	for range tracker.cleaner.C {
		tracker.mutex.Lock()
		for fileHash := range tracker.index {
			for peerAddress, peerTimestamp := range tracker.index[fileHash] {
				if time.Since(peerTimestamp) > tracker.config.AnnounceInterval*2 {
					log.Println("clean", fileHash, peerAddress)
					delete(tracker.index[fileHash], peerAddress)
				}
				if len(tracker.index[fileHash]) == 0 {
					log.Println("clean", fileHash)
					delete(tracker.index, fileHash)
				}
			}
		}
		tracker.mutex.Unlock()
	}
}

func (tracker *Tracker) Announce(ctx context.Context, req *AnnounceRequest) (*AnnounceResponse, error) {
	log.Println("Announce", req.FileHash, req.PeerAddress)

	tracker.mutex.Lock()
	if _, ok := tracker.index[req.FileHash]; !ok {
		tracker.index[req.FileHash] = make(map[string]time.Time)
	}
	tracker.mutex.Unlock()

	tracker.mutex.Lock()
	tracker.index[req.FileHash][req.PeerAddress] = time.Now()
	tracker.mutex.Unlock()

	tracker.mutex.RLock()
	i := 0
	peerAddresses := make([]string, len(tracker.index[req.FileHash])-1)
	for peerAddress := range tracker.index[req.FileHash] {
		// Don't return announcing peer his own address
		if peerAddress == req.PeerAddress {
			continue
		}

		peerAddresses[i] = peerAddress
		i++
	}
	tracker.mutex.RUnlock()

	return &AnnounceResponse{
		PeerAddresses:    peerAddresses,
		AnnounceInterval: tracker.config.AnnounceInterval.Milliseconds() / 1000,
	}, nil
}
