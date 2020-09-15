package main

import (
	"context"
	"log"
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
	log.Println("Add")
	return &AddResponse{Torrent: in.Torrent}, nil
}
