// Package tracker implements tracker service
package tracker

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/dawidd6/p2p/pkg/config"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	PeerFromContextError = errors.New("can't determine peer from context")
)

// Tracker represents tracker service
type Tracker struct {
	conf *config.Config

	index      map[string]map[string]time.Time // fileHash peerAddress peerTimestamp
	indexMutex sync.RWMutex

	cleanTicker *time.Ticker

	UnimplementedTrackerServer
}

// Run starts tracker server
func Run(conf *config.Config) error {
	tracker := &Tracker{
		conf:        conf,
		index:       make(map[string]map[string]time.Time),
		cleanTicker: time.NewTicker(conf.CleanInterval),
	}

	// Start cleaning peer index
	go tracker.cleaning()

	// Listen on address
	address := net.JoinHostPort(conf.TrackerHost, conf.TrackerPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// Register and serve
	server := grpc.NewServer()
	RegisterTrackerServer(server, tracker)
	return server.Serve(listener)
}

// cleaning prunes peers periodically
func (tracker *Tracker) cleaning() {
	for range tracker.cleanTicker.C {
		tracker.indexMutex.Lock()
		tracker.clean()
		tracker.indexMutex.Unlock()
	}
}

// clean prunes dangling peers
func (tracker *Tracker) clean() {
	now := time.Now()

	for fileHash, peerInfo := range tracker.index {
		for peerAddress, peerTimestamp := range peerInfo {
			// Delete peer entry
			if now.Sub(peerTimestamp) > tracker.conf.AnnounceInterval*2 {
				delete(tracker.index[fileHash], peerAddress)
			}
			// Delete torrent entry
			if len(tracker.index[fileHash]) == 0 {
				delete(tracker.index, fileHash)
			}
		}
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
	announceInterval := tracker.conf.AnnounceInterval.Milliseconds() / 1000

	// Create torrent file hash entry in map, if does not exist
	tracker.indexMutex.Lock()
	if _, ok := tracker.index[req.FileHash]; !ok {
		tracker.index[req.FileHash] = make(map[string]time.Time)
	}
	tracker.indexMutex.Unlock()

	// Add peer address into map
	tracker.indexMutex.Lock()
	tracker.index[req.FileHash][peerAddress] = time.Now()
	tracker.indexMutex.Unlock()

	// Get list of peers for a torrent
	tracker.indexMutex.RLock()
	peerAddresses := make([]string, 0, len(tracker.index[req.FileHash])-1)
	for peerAddr := range tracker.index[req.FileHash] {
		// Don't return announcing peer their own address
		if peerAddr != peerAddress {
			peerAddresses = append(peerAddresses, peerAddr)
		}
	}
	tracker.indexMutex.RUnlock()

	return &AnnounceResponse{
		PeerAddresses:    peerAddresses,
		AnnounceInterval: announceInterval,
	}, nil
}
