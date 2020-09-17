package daemon

import (
	"context"
	"log"

	"github.com/dawidd6/p2p/torrent"
	"github.com/dawidd6/p2p/tracker"

	"google.golang.org/grpc"
)

type Daemon struct {
	torrents []*torrent.Torrent
}

func NewDaemon() *Daemon {
	return &Daemon{
		torrents: make([]*torrent.Torrent, 0),
	}
}

func (daemon *Daemon) Add(ctx context.Context, in *AddRequest) (*AddResponse, error) {
	for _, url := range in.Torrent.Trackers {
		conn, err := grpc.Dial(url, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		request := &tracker.AnnounceRequest{
			TorrentSha256: in.Torrent.Sha256,
		}
		client := tracker.NewTrackerClient(conn)
		_, err = client.Announce(context.TODO(), request)
		if err != nil {
			return nil, err
		}
	}

	log.Println("Add")

	return &AddResponse{Torrent: in.Torrent}, nil
}
