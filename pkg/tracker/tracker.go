// Package tracker implements tracker service
package tracker

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/dawidd6/p2p/pkg/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	PeerFromContextError = errors.New("can't determine peer from context")
)

// Tracker represents tracker service
type Tracker struct {
	conf *config.Config
	db   *redis.Pool

	UnimplementedTrackerServer
}

// Run starts tracker server
func Run(conf *config.Config) error {
	// Create connection pool to database
	db := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			address := net.JoinHostPort(conf.DBHost, conf.DBPort)
			return redis.Dial("tcp", address)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Wait: true,
	}

	// Construct tracker
	tracker := &Tracker{
		conf: conf,
		db:   db,
	}

	channel := make(chan error)

	// Keep testing database connection and availability
	go func() {
		for {
			conn := db.Get()
			err := conn.Err()
			if err != nil {
				channel <- err
			}
			defer conn.Close()

			time.Sleep(conf.DBCheckInterval)
		}
	}()

	// Start tracker server
	go func() {
		address := net.JoinHostPort(conf.TrackerHost, conf.TrackerPort)
		listener, err := net.Listen("tcp", address)
		if err != nil {
			channel <- err
		}

		log.Println("Listening on", address)

		server := grpc.NewServer()
		RegisterTrackerServer(server, tracker)
		channel <- server.Serve(listener)
	}()

	return <-channel
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
	peerCutoff := (tracker.conf.AnnounceInterval + tracker.conf.AnnounceTolerance).Milliseconds() / 1000

	log.Println("Announcing", peerAddress, "for", req.FileHash)

	// Get database connection from pool
	conn, err := tracker.db.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	// Close connection to database on exit from this function
	defer conn.Close()

	// Add peer address to set
	_, err = conn.Do("SADD", req.FileHash, peerAddress)
	if err != nil {
		return nil, err
	}

	// Set peer address expire time
	_, err = conn.Do("EXPIREMEMBER", req.FileHash, peerAddress, peerCutoff)
	if err != nil {
		return nil, err
	}

	// Get peer addresses for given file hash
	peerAddresses, err := redis.Strings(conn.Do("SMEMBERS", req.FileHash))
	if err != nil {
		return nil, err
	}

	// Remove announcing peer's address from the list
	for i := range peerAddresses {
		if peerAddresses[i] == peerAddress {
			peerAddresses = append(peerAddresses[:i], peerAddresses[i+1:]...)
			break
		}
	}

	// Return to client
	return &AnnounceResponse{
		PeerAddresses:    peerAddresses,
		AnnounceInterval: announceInterval,
	}, nil
}
