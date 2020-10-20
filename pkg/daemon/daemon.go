package daemon

import (
	"context"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/dawidd6/p2p/pkg/hash"
	"github.com/dawidd6/p2p/pkg/piece"

	"github.com/dawidd6/p2p/pkg/file"

	"github.com/dawidd6/p2p/pkg/errors"
	"github.com/dawidd6/p2p/pkg/torrent"
	"github.com/dawidd6/p2p/pkg/tracker"
	"google.golang.org/grpc"
)

type Config struct {
	DownloadsDir      string
	ListenAddress     string
	SeedListenAddress string
}

type State struct {
	DownloadedPieces uint64
	File             *os.File
	Torrent          *torrent.Torrent
	Peers            []string
	AnnounceInterval time.Duration
	AnnounceChannel  chan struct{}
}

type Daemon struct {
	config      *Config
	torrents    map[string]*State
	mutex       sync.RWMutex
	addChannel  chan *AddRequest
	seedChannel chan *SeedRequest
	UnimplementedDaemonServer
	UnimplementedSeederServer
}

func Run(config *Config) error {
	daemon := &Daemon{
		config:      config,
		torrents:    make(map[string]*State),
		addChannel:  make(chan *AddRequest),
		seedChannel: make(chan *SeedRequest, 5),
	}

	channel := make(chan error)

	go func() {
		listener, err := net.Listen("tcp", config.ListenAddress)
		if err != nil {
			channel <- err
		}

		server := grpc.NewServer()
		RegisterDaemonServer(server, daemon)
		channel <- server.Serve(listener)
	}()

	go func() {
		listener, err := net.Listen("tcp", config.SeedListenAddress)
		if err != nil {
			channel <- err
		}

		server := grpc.NewServer()
		RegisterSeederServer(server, daemon)
		channel <- server.Serve(listener)
	}()

	// channel <- daemon.loop() TODO

	return <-channel
}

func (daemon *Daemon) announce(state *State) error {
	conn, err := grpc.Dial(state.Torrent.TrackerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}

	request := &tracker.AnnounceRequest{
		FileHash:    state.Torrent.FileHash,
		PeerAddress: daemon.config.SeedListenAddress,
	}

	response, err := tracker.NewTrackerClient(conn).Announce(context.TODO(), request)
	if err != nil {
		return err
	}

	state.Peers = response.PeerAddresses
	state.AnnounceInterval = time.Duration(response.AnnounceInterval) * time.Second

	return nil
}

func (daemon *Daemon) fetch(state *State) {
	err := file.CreateDirectory(daemon.config.DownloadsDir)

	if err != nil {
		log.Println(err)
		return
	}

	err = daemon.announce(state)
	if err != nil {
		log.Println(err)
	}

	go func() {
		for {
			select {
			case <-time.After(state.AnnounceInterval):
				err = daemon.announce(state)
			case <-state.AnnounceChannel:
				break
			}
		}
	}()

	state.File, err = file.Open(daemon.config.DownloadsDir, state.Torrent.FileName)
	if err != nil {
		log.Println(err)
		return
	}

	err = torrent.Verify(state.Torrent, daemon.config.DownloadsDir)
	if err == nil {
		log.Println("seeding")
		return
	}

	for {
		if len(state.Peers) == 0 {
			log.Println("no peers available")
			<-time.After(time.Second * 5)
			continue
		}

		for i, pieceHash := range state.Torrent.PieceHashes {
			pieceNumber := uint64(i)

			// TODO, throttle for now
			<-time.After(time.Second * 1)

			request := &SeedRequest{
				FileHash:    state.Torrent.FileHash,
				PieceNumber: pieceNumber,
			}

			pieceData, err := piece.Read(state.File, state.Torrent.PieceSize, pieceNumber)
			if err != nil {
				log.Println(err)
			}

			if hash.Compute(pieceData) == pieceHash {
				log.Println(pieceHash, "already downloaded piece")
				continue
			}

			for _, peerAddr := range state.Peers {
				log.Println("Fetch", pieceNumber, pieceHash)

				// TODO use DialContext everywhere
				conn, err := grpc.DialContext(context.TODO(), peerAddr, grpc.WithInsecure())
				if err != nil {
					log.Println(err)
					continue
				}

				response, err := NewSeederClient(conn).Seed(context.TODO(), request)
				if err != nil {
					log.Println(err)
					continue
				}

				// Peer does not have this piece, ask someone else
				if response.PieceData == nil {
					log.Println(peerAddr, "no piece")
					continue
				}

				// Got the piece wrongly, ask someone else
				if hash.Compute(response.PieceData) != pieceHash {
					log.Println(peerAddr, "wrong piece")
					continue
				}

				err = piece.Write(state.File, pieceNumber, pieceData)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}

		err = torrent.Verify(state.Torrent, daemon.config.DownloadsDir)
		if err != nil {
			log.Println(err)
			continue
		} else {
			log.Println("completed")
			return
		}
	}
}

func (daemon *Daemon) Seed(ctx context.Context, req *SeedRequest) (*SeedResponse, error) {
	daemon.mutex.RLock()
	defer daemon.mutex.RUnlock()

	state, ok := daemon.torrents[req.FileHash]
	if !ok {
		return nil, errors.TorrentNotFound
	}

	pieceData, err := piece.Read(state.File, state.Torrent.PieceSize, req.PieceNumber)
	if err != nil {
		return nil, err
	}

	return &SeedResponse{PieceData: pieceData}, nil
}

func (daemon *Daemon) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	daemon.mutex.Lock()
	defer daemon.mutex.Unlock()

	state, ok := daemon.torrents[req.Torrent.FileHash]
	if ok {
		return nil, errors.TorrentAlreadyAdded
	}

	state = &State{
		DownloadedPieces: 0,
		Torrent:          req.Torrent,
		PeerAddresses:    []string{},
		AnnounceChannel:  make(chan struct{}),
		AnnounceInterval: defaults.TrackerAnnounceInterval,
	}

	daemon.torrents[req.Torrent.FileHash] = state
	go daemon.fetch(state)

	return &AddResponse{Torrent: req.Torrent}, nil
}
