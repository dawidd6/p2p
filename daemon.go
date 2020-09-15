package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

type Daemon struct {
	torrents []*Torrent
}

func NewDaemon() *Daemon {
	return &Daemon{
		torrents: make([]*Torrent, 0),
	}
}

func (daemon *Daemon) Add(ctx context.Context, in *AddRequest) (*AddResponse, error) {
	for _, url := range in.Torrent.Trackers {
		conn, err := grpc.Dial(url, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		request := &AnnounceRequest{
			TorrentSha256: in.Torrent.Sha256,
		}
		tracker := NewTrackerClient(conn)
		_, err = tracker.Announce(context.TODO(), request)
		if err != nil {
			return nil, err
		}
	}

	log.Println("Add")

	return &AddResponse{Torrent: in.Torrent}, nil
}
