// Package tracker implements tracker service
package tracker

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const (
	Host             = "127.0.0.1"
	Port             = "8889"
	AnnounceInterval = time.Second * 30
	CleanInterval    = time.Second * 60
	ConnTimeout      = time.Second * 5
)

var (
	PeerFromContextError = errors.New("can't determine peer from context")
)

// Config holds all configurable aspects of tracker
type Config struct {
	Host             string
	Port             string
	AnnounceInterval time.Duration
	CleanInterval    time.Duration
}

// Tracker represents tracker service
type Tracker struct {
	config  *Config
	index   map[string]map[string]time.Time // fileHash peerAddress peerTimestamp
	mutex   sync.RWMutex
	cleaner *time.Ticker
	UnimplementedTrackerServer
}

// Run starts tracker server
func Run(config *Config) error {
	tracker := &Tracker{
		config:  config,
		index:   make(map[string]map[string]time.Time),
		cleaner: time.NewTicker(config.CleanInterval),
	}

	// Start cleaning peer index
	go tracker.clean()

	// Listen on address
	address := net.JoinHostPort(config.Host, config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// Register and serve
	server := grpc.NewServer(grpc.ConnectionTimeout(ConnTimeout))
	RegisterTrackerServer(server, tracker)
	return server.Serve(listener)
}

// clean prunes peers periodically
func (tracker *Tracker) clean() {
	for range tracker.cleaner.C {
		tracker.mutex.Lock()
		for fileHash := range tracker.index {
			for peerAddress, peerTimestamp := range tracker.index[fileHash] {
				// Delete peer entry
				if time.Since(peerTimestamp) > tracker.config.AnnounceInterval*2 {
					log.Println("clean", fileHash, peerAddress)
					delete(tracker.index[fileHash], peerAddress)
				}
				// Delete torrent entry
				if len(tracker.index[fileHash]) == 0 {
					log.Println("clean", fileHash)
					delete(tracker.index, fileHash)
				}
			}
		}
		tracker.mutex.Unlock()
	}
}

// Announce is called by the daemon and it adds the peer to index
func (tracker *Tracker) Announce(ctx context.Context, req *AnnounceRequest) (*AnnounceResponse, error) {
	// Get peer info from context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, PeerFromContextError
	}

	// Split host:port peer address string
	host, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return nil, err
	}

	// Construct needed variables
	peerAddress := net.JoinHostPort(host, req.PeerPort)
	announceInterval := tracker.config.AnnounceInterval.Milliseconds() / 1000

	log.Println("Announce", req.FileHash, peerAddress)

	// Create torrent file hash entry in map, if does not exist
	tracker.mutex.Lock()
	if _, ok := tracker.index[req.FileHash]; !ok {
		tracker.index[req.FileHash] = make(map[string]time.Time)
	}
	tracker.mutex.Unlock()

	// Add peer address into map
	tracker.mutex.Lock()
	tracker.index[req.FileHash][peerAddress] = time.Now()
	tracker.mutex.Unlock()

	// Get list of peers for a torrent
	tracker.mutex.RLock()
	peerAddresses := make([]string, 0, len(tracker.index[req.FileHash])-1)
	for peerAddressKey := range tracker.index[req.FileHash] {
		// Don't return announcing peer his own address
		if peerAddressKey != peerAddress {
			peerAddresses = append(peerAddresses, peerAddressKey)
		}
	}
	tracker.mutex.RUnlock()

	return &AnnounceResponse{
		PeerAddresses:    peerAddresses,
		AnnounceInterval: announceInterval,
	}, nil
}
