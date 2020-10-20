package daemon

import (
	"context"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/dawidd6/p2p/pkg/state"

	"github.com/dawidd6/p2p/pkg/defaults"
	"github.com/dawidd6/p2p/pkg/hash"
	"github.com/dawidd6/p2p/pkg/piece"

	"github.com/dawidd6/p2p/pkg/errors"
	"github.com/dawidd6/p2p/pkg/torrent"
	"github.com/dawidd6/p2p/pkg/tracker"
	"google.golang.org/grpc"
)

type Config struct {
	MaxWorkers        int
	DownloadsDir      string
	ListenAddress     string
	SeedListenAddress string
}

type Task struct {
	File             *os.File
	Torrent          *torrent.Torrent
	State            *state.State
	PeerAddresses    []string
	AnnounceInterval time.Duration
}

type Daemon struct {
	config   *Config
	torrents map[string]*Task
	mutex    sync.RWMutex
	UnimplementedDaemonServer
	UnimplementedSeederServer
}

func Run(config *Config) error {
	daemon := &Daemon{
		config:   config,
		torrents: make(map[string]*Task),
	}

	err := os.MkdirAll(config.DownloadsDir, 0775)
	if err != nil {
		return err
	}

	err = os.Chdir(config.DownloadsDir)
	if err != nil {
		return err
	}

	channel := make(chan error)

	go func() {
		listener, err := net.Listen("tcp", config.ListenAddress)
		if err != nil {
			channel <- err
		}

		server := grpc.NewServer(grpc.ConnectionTimeout(defaults.DaemonConnTimeout))
		RegisterDaemonServer(server, daemon)
		channel <- server.Serve(listener)
	}()

	go func() {
		listener, err := net.Listen("tcp", config.SeedListenAddress)
		if err != nil {
			channel <- err
		}

		server := grpc.NewServer(grpc.ConnectionTimeout(defaults.SeedConnTimeout))
		RegisterSeederServer(server, daemon)
		channel <- server.Serve(listener)
	}()

	// channel <- daemon.loop() TODO

	return <-channel
}

func (daemon *Daemon) announce(task *Task) error {
	conn, err := grpc.Dial(task.Torrent.TrackerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}

	request := &tracker.AnnounceRequest{
		FileHash:    task.Torrent.FileHash,
		PeerAddress: daemon.config.SeedListenAddress,
	}

	client := tracker.NewTrackerClient(conn)
	response, err := client.Announce(context.Background(), request)
	if err != nil {
		return err
	}

	err = conn.Close()
	if err != nil {
		return err
	}

	task.PeerAddresses = response.PeerAddresses
	task.AnnounceInterval = time.Duration(response.AnnounceInterval) * time.Second

	return nil
}

func (daemon *Daemon) fetch(task *Task) {
	err := daemon.announce(task)
	if err != nil {
		log.Println(err)
	}

	go func() {
		for {
			// TODO break if deleted torrent
			select {
			case <-time.After(task.AnnounceInterval):
				err = daemon.announce(task)
			}
		}
	}()

	task.File, err = os.OpenFile(task.Torrent.FileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	err = torrent.Verify(task.Torrent, task.File)
	if err == nil {
		log.Println("seeding")
		return
	}

	for {
		if len(task.PeerAddresses) == 0 {
			log.Println("no peers available")
			<-time.After(task.AnnounceInterval)
			continue
		}

		wg := sync.WaitGroup{}
		workers := make(chan struct{}, daemon.config.MaxWorkers)

		for i := range task.Torrent.PieceHashes {
			pieceNumber := int64(i)
			pieceHash := task.Torrent.PieceHashes[pieceNumber]
			pieceOffset := piece.Offset(task.Torrent.PieceSize, pieceNumber)

			workers <- struct{}{}
			wg.Add(1)

			go func() {
				defer func() {
					<-workers
					wg.Done()
				}()

				pieceData, err := piece.Read(task.File, task.Torrent.PieceSize, pieceOffset)
				if err != nil {
					log.Println(err)
				}

				err = hash.New().Verify(pieceData, pieceHash)
				if err == nil {
					log.Println(pieceHash, "already downloaded piece")
					return
				}

				request := &SeedRequest{
					FileHash:    task.Torrent.FileHash,
					PieceNumber: pieceNumber,
				}

				for _, peerAddr := range task.PeerAddresses {
					log.Println("Fetch", pieceNumber, pieceHash, peerAddr)

					conn, err := grpc.Dial(peerAddr, grpc.WithInsecure())
					if err != nil {
						log.Println(err)
						continue
					}

					client := NewSeederClient(conn)
					response, err := client.Seed(context.Background(), request)
					if err != nil {
						log.Println(err)
						continue
					}

					err = conn.Close()
					if err != nil {
						log.Println(err)
						continue
					}

					// Got the piece wrongly, ask someone else
					err = hash.New().Verify(response.PieceData, pieceHash)
					if err != nil {
						log.Println(peerAddr, "wrong piece")
						continue
					}

					err = piece.Write(task.File, pieceOffset, response.PieceData)
					if err != nil {
						log.Println(err)
						continue
					}

					break
				}
			}()
		}

		wg.Wait()

		err = torrent.Verify(task.Torrent, task.File)
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
	task, ok := daemon.torrents[req.FileHash]
	daemon.mutex.RUnlock()
	if !ok {
		return nil, errors.TorrentNotFound
	}

	pieceOffset := piece.Offset(task.Torrent.PieceSize, req.PieceNumber)
	pieceData, err := piece.Read(task.File, task.Torrent.PieceSize, pieceOffset)
	if err != nil {
		return nil, err
	}

	return &SeedResponse{PieceData: pieceData}, nil
}

func (daemon *Daemon) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	daemon.mutex.RLock()
	task, ok := daemon.torrents[req.Torrent.FileHash]
	daemon.mutex.RUnlock()
	if ok {
		return nil, errors.TorrentAlreadyAdded
	}

	task = &Task{
		Torrent:          req.Torrent,
		AnnounceInterval: defaults.TrackerAnnounceInterval,
	}

	daemon.mutex.Lock()
	daemon.torrents[req.Torrent.FileHash] = task
	daemon.mutex.Unlock()

	go daemon.fetch(task)

	return &AddResponse{}, nil
}

func (daemon *Daemon) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	daemon.mutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.mutex.RUnlock()
	if !ok {
		return nil, errors.TorrentNotFound
	}

	daemon.mutex.Lock()
	delete(daemon.torrents, req.FileHash)
	daemon.mutex.Unlock()

	err := task.File.Close()
	if err != nil {
		return nil, err
	}

	if req.WithData {
		err := os.Remove(task.Torrent.FileName)
		if err != nil {
			return nil, err
		}
	}

	return &DeleteResponse{}, nil
}

func (daemon *Daemon) Status(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	daemon.mutex.RLock()
	torrents := make(map[string]*state.State, len(daemon.torrents))
	for fileHash := range daemon.torrents {
		torrents[fileHash] = daemon.torrents[fileHash].State
	}
	daemon.mutex.RUnlock()

	return &StatusResponse{
		Torrents: torrents,
	}, nil
}
