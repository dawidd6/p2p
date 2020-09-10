package main

import (
	"context"
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
	return &AddResponse{Torrent: in.Torrent}, nil
}

func (daemon *Daemon) Delete(ctx context.Context, in *DeleteRequest) (*DeleteResponse, error) {
	return &DeleteResponse{Torrent: in.Torrent}, nil
}
