package daemon

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/dawidd6/p2p/pkg/worker"

	"github.com/dawidd6/p2p/pkg/state"

	"github.com/dawidd6/p2p/pkg/hash"
	"github.com/dawidd6/p2p/pkg/piece"

	"github.com/dawidd6/p2p/pkg/torrent"
	"github.com/dawidd6/p2p/pkg/tracker"
	"google.golang.org/grpc"
)

const (
	ListenAddress      = "127.0.0.1:8888"
	SeedListenAddress  = "0.0.0.0:44444"
	ConnTimeout        = time.Second * 3
	SeedConnTimeout    = time.Second * 2
	DownloadsDir       = "."
	MaxFetchingWorkers = 4
)

var (
	TorrentNotFoundError     = errors.New("torrent not found")
	TorrentAlreadyAddedError = errors.New("torrent is already added")
)

type Config struct {
	MaxFetchingWorkers int
	DownloadsDir       string
	ListenAddress      string
	SeedListenAddress  string
}

type Task struct {
	File               *os.File
	Torrent            *torrent.Torrent
	State              *state.State
	PeerAddresses      []string
	PeerAddressesMutex sync.RWMutex
	PeersAvailable     sync.WaitGroup
	AnnounceInterval   time.Duration
	AnnounceTicker     *time.Ticker
	Randomizer         *rand.Rand
	Resume             sync.WaitGroup
	WorkerPool         *worker.Pool
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

		server := grpc.NewServer(
			grpc.ConnectionTimeout(ConnTimeout),
		)
		RegisterDaemonServer(server, daemon)
		channel <- server.Serve(listener)
	}()

	go func() {
		listener, err := net.Listen("tcp", config.SeedListenAddress)
		if err != nil {
			channel <- err
		}

		server := grpc.NewServer(
			grpc.ConnectionTimeout(SeedConnTimeout),
		)
		RegisterSeederServer(server, daemon)
		channel <- server.Serve(listener)
	}()

	return <-channel
}

func (daemon *Daemon) announce(task *Task) {
	for range task.AnnounceTicker.C {
		conn, err := grpc.Dial(task.Torrent.TrackerAddress, grpc.WithInsecure())
		if err != nil {
			log.Println("announce", err)
		}

		request := &tracker.AnnounceRequest{
			FileHash:    task.Torrent.FileHash,
			PeerAddress: daemon.config.SeedListenAddress,
		}

		client := tracker.NewTrackerClient(conn)
		response, err := client.Announce(context.Background(), request)
		if err != nil {
			log.Println("announce", err)
		}

		err = conn.Close()
		if err != nil {
			log.Println("announce", err)
		}

		task.PeerAddressesMutex.Lock()
		task.PeerAddresses = response.PeerAddresses
		task.PeerAddressesMutex.Unlock()

		announceInterval := time.Duration(response.AnnounceInterval) * time.Second
		if task.AnnounceInterval != announceInterval {
			log.Println("announce", announceInterval)
			task.AnnounceInterval = announceInterval
			task.AnnounceTicker.Reset(task.AnnounceInterval)
		}
	}
}

func (daemon *Daemon) fetch(task *Task, pieceNumber int64, pieceHash string, pieceOffset int64) {
	hasher := hash.New()

	// Try reading piece from disk
	pieceData, err := piece.Read(task.File, task.Torrent.PieceSize, pieceOffset)
	if err != nil {
		log.Println(err)
	}

	// Check if read piece is correct
	err = hasher.Verify(pieceData, pieceHash)
	if err == nil {
		log.Println(pieceHash, "already downloaded piece")
		return
	}

	// Construct seed request
	request := &SeedRequest{
		FileHash:    task.Torrent.FileHash,
		PieceNumber: pieceNumber,
	}

	// Get random peer address
	task.PeerAddressesMutex.RLock()
	i := task.Randomizer.Intn(len(task.PeerAddresses))
	peerAddr := task.PeerAddresses[i]
	task.PeerAddressesMutex.RUnlock()

	// Connect to peer
	conn, err := grpc.Dial(peerAddr, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return
	}

	// Get the piece from peer
	client := NewSeederClient(conn)
	response, err := client.Seed(context.Background(), request)
	if err != nil {
		log.Println(err)
		return
	}

	// Close the connection
	err = conn.Close()
	if err != nil {
		log.Println(err)
		return
	}

	// Got the piece wrongly, ask someone else
	err = hasher.Verify(response.PieceData, pieceHash)
	if err != nil {
		log.Println(peerAddr, "wrong piece")
		return
	}

	// Save downloaded piece on disk, can be called concurrently
	err = piece.Write(task.File, pieceOffset, response.PieceData)
	if err != nil {
		log.Println(err)
		return
	}
}

func (daemon *Daemon) add(task *Task) {
	var err error

	log.Println("add")

	// Open or create the torrent data file
	task.File, err = os.OpenFile(task.Torrent.FileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	// Check if torrent is already completed, if it is - just start seeding
	err = torrent.Verify(task.Torrent, task.File)
	if err == nil {
		log.Println("seeding")
		return
	}

	// Start a fetch loop
	for {
		// start a worker pool
		task.WorkerPool.Start()

		// Loop over torrent piece hashes
		for i := range task.Torrent.PieceHashes {
			pieceNumber := int64(i)
			pieceHash := task.Torrent.PieceHashes[pieceNumber]
			pieceOffset := piece.Offset(task.Torrent.PieceSize, pieceNumber)

			// Wait here if torrent is paused, will continue if it is resumed
			task.Resume.Wait()

			// Create worker, wait if max count already
			task.WorkerPool.Enqueue(func() {
				// Fetch one piece
				daemon.fetch(task, pieceNumber, pieceHash, pieceOffset)
			})
		}

		// Wait for all fetch workers to complete
		task.WorkerPool.Finish()

		// Verify if all downloaded pieces are correct
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
	// Check if torrent is present in queue
	daemon.mutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.mutex.RUnlock()
	if !ok {
		return nil, TorrentNotFoundError
	}

	// Wait if torrent is paused
	task.Resume.Wait()

	// Read piece
	pieceOffset := piece.Offset(task.Torrent.PieceSize, req.PieceNumber)
	pieceData, err := piece.Read(task.File, task.Torrent.PieceSize, pieceOffset)
	if err != nil {
		return nil, err
	}

	// Return to client
	return &SeedResponse{PieceData: pieceData}, nil
}

func (daemon *Daemon) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	// Check if torrent is present in queue
	daemon.mutex.RLock()
	task, ok := daemon.torrents[req.Torrent.FileHash]
	daemon.mutex.RUnlock()
	if ok {
		return nil, TorrentAlreadyAddedError
	}

	// Create new task
	task = &Task{
		Torrent:          req.Torrent,
		AnnounceInterval: tracker.AnnounceInterval,
		AnnounceTicker:   time.NewTicker(tracker.AnnounceInterval),
		Randomizer:       rand.New(rand.NewSource(time.Now().Unix())),
		WorkerPool:       worker.New(daemon.config.MaxFetchingWorkers),
	}

	// Assign task to torrent
	daemon.mutex.Lock()
	daemon.torrents[req.Torrent.FileHash] = task
	daemon.mutex.Unlock()

	// Keep announcing torrent to tracker
	go daemon.announce(task)
	// Start torrent process
	go daemon.add(task)

	// Return to client
	return &AddResponse{}, nil
}

func (daemon *Daemon) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	// Check if torrent is present in queue
	daemon.mutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.mutex.RUnlock()
	if !ok {
		return nil, TorrentNotFoundError
	}

	// Delete torrent from queue
	daemon.mutex.Lock()
	delete(daemon.torrents, req.FileHash)
	daemon.mutex.Unlock()

	// Stop announcing torrent to tracker
	task.AnnounceTicker.Stop()

	// Finish any fetching workers
	task.WorkerPool.Finish()

	// Close torrent file
	err := task.File.Close()
	if err != nil {
		return nil, err
	}

	// Remove downloaded torrent data if desired
	if req.WithData {
		err := os.Remove(task.Torrent.FileName)
		if err != nil {
			return nil, err
		}
	}

	// Return to client
	return &DeleteResponse{}, nil
}

func (daemon *Daemon) Status(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	// Handle case when status is requested for just one torrent
	if req.FileHash != "" {
		// Check if torrent is present in queue
		daemon.mutex.RLock()
		task, ok := daemon.torrents[req.FileHash]
		daemon.mutex.RUnlock()
		if !ok {
			return nil, TorrentNotFoundError
		}

		// Return to client
		return &StatusResponse{
			Torrents: map[string]*state.State{
				task.Torrent.FileHash: task.State,
			},
		}, nil
	}

	// Get all torrents and construct a status map
	daemon.mutex.RLock()
	torrents := make(map[string]*state.State, len(daemon.torrents))
	for fileHash := range daemon.torrents {
		torrents[fileHash] = daemon.torrents[fileHash].State
	}
	daemon.mutex.RUnlock()

	// Return to client
	return &StatusResponse{
		Torrents: torrents,
	}, nil
}

func (daemon *Daemon) Resume(ctx context.Context, req *ResumeRequest) (*ResumeResponse, error) {
	// Check if torrent is present in queue
	daemon.mutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.mutex.RUnlock()
	if !ok {
		return nil, TorrentNotFoundError
	}

	// Resume torrent
	task.Resume.Done()

	// Return to client
	return &ResumeResponse{}, nil
}

func (daemon *Daemon) Pause(ctx context.Context, req *PauseRequest) (*PauseResponse, error) {
	// Check if torrent is present in queue
	daemon.mutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.mutex.RUnlock()
	if !ok {
		return nil, TorrentNotFoundError
	}

	// Pause torrent
	task.Resume.Add(1)

	// Return to client
	return &PauseResponse{}, nil
}
