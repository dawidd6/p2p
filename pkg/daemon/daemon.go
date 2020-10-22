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
	TorrentNotFoundError       = errors.New("torrent not found")
	TorrentAlreadyAddedError   = errors.New("torrent is already added")
	TorrentAlreadyPausedError  = errors.New("torrent is already paused")
	TorrentAlreadyResumedError = errors.New("torrent is already resumed")
)

type Config struct {
	MaxFetchingWorkers int
	DownloadsDir       string
	ListenAddress      string
	SeedListenAddress  string
}

type Task struct {
	File    *os.File
	Torrent *torrent.Torrent
	State   *state.State

	ResumeNotifier chan struct{}
	PauseNotifier  chan struct{}
	DeleteNotifier chan struct{}

	PeersNotifier   chan struct{}
	PeersMutex      sync.Mutex
	Peers           []string
	PeersRandomizer *rand.Rand

	AnnounceTicker   *time.Ticker
	AnnounceInterval time.Duration

	WorkerPool *worker.Pool
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

func (daemon *Daemon) announce(task *Task) error {
	log.Println("announce", "start")
	defer log.Println("announce", "stop")

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

	// Set peer addresses and notify about that fact
	task.PeersMutex.Lock()
	task.Peers = response.PeerAddresses
	if len(task.PeersNotifier) == 1 {
		log.Println("announce", "draining channel")
		<-task.PeersNotifier
		log.Println("announce", "channel drained")
	}
	log.Println("announce", "notifying")
	task.PeersNotifier <- struct{}{}
	log.Println("announce", "notified")
	task.PeersMutex.Unlock()

	// Set announcing interval if changed and reset the ticker
	announceInterval := time.Duration(response.AnnounceInterval) * time.Second
	if task.AnnounceInterval != announceInterval {
		log.Println("announce", "setting interval", announceInterval)
		task.AnnounceInterval = announceInterval
		task.AnnounceTicker.Reset(task.AnnounceInterval)
	}

	return nil
}

func (daemon *Daemon) announcing(task *Task) {
	for range task.AnnounceTicker.C {
		err := daemon.announce(task)
		if err != nil {
			log.Println(err)
		}
	}
}

func (daemon *Daemon) fetch(task *Task, pieceNumber int64, pieceHash string, pieceOffset int64, peerAddr string) error {
	log.Println("fetch", pieceNumber, peerAddr)

	hasher := hash.New()

	// Try reading piece from disk
	pieceData, err := piece.Read(task.File, task.Torrent.PieceSize, pieceOffset)
	if err != nil {
		return err
	}

	// Check if read piece is correct, return if it is
	err = hasher.Verify(pieceData, pieceHash)
	if err == nil {
		return nil
	}

	// Construct seed request
	request := &SeedRequest{
		FileHash:    task.Torrent.FileHash,
		PieceNumber: pieceNumber,
	}

	// Connect to peer
	conn, err := grpc.Dial(peerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	// Get the piece from peer
	client := NewSeederClient(conn)
	response, err := client.Seed(context.Background(), request)
	if err != nil {
		return err
	}

	// Close the connection
	err = conn.Close()
	if err != nil {
		return err
	}

	// Got the piece wrongly, ask someone else
	err = hasher.Verify(response.PieceData, pieceHash)
	if err != nil {
		return err
	}

	// Save downloaded piece on disk, can be called concurrently
	err = piece.Write(task.File, pieceOffset, response.PieceData)
	if err != nil {
		return err
	}

	task.State.DownloadedBytes += int64(len(response.PieceData))

	return nil
}

func (daemon *Daemon) fetching(task *Task) {
	log.Println("add")

	err := daemon.announce(task)
	if err != nil {
		log.Println(err)
	}

	// Open or create the torrent data file
	task.File, err = os.OpenFile(task.Torrent.FileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	// Check if torrent is already completed, if it is - just start seeding
	err = torrent.Verify(task.Torrent, task.File)
	if err == nil {
		task.State.DownloadedBytes = task.Torrent.FileSize
		task.State.Completed = true
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

			// Pause execution if desired or exit from function if torrent is deleted
			select {
			case <-task.PauseNotifier:
				log.Println("fetch", "paused", task.Torrent.FileHash)
				<-task.ResumeNotifier
				log.Println("fetch", "resumed", task.Torrent.FileHash)
			case <-task.DeleteNotifier:
				log.Println("fetch", "deleted", task.Torrent.FileHash)
				return
			default:
			}

			// Get random peer address, wait if no peers available
			peerAddr := ""
			for {
				task.PeersMutex.Lock()
				peers := len(task.Peers)
				if peers > 0 {
					peer := task.PeersRandomizer.Intn(peers)
					peerAddr = task.Peers[peer]
					task.PeersMutex.Unlock()
					log.Println("fetch", "peer", peerAddr)
					break
				}
				task.PeersMutex.Unlock()
				log.Println("fetch", "waiting for peers")
				<-task.PeersNotifier
				log.Println("fetch", "got new peers")
			}

			// Create worker, wait if max count already
			task.WorkerPool.Enqueue(func() {
				// Fetch one piece
				err := daemon.fetch(task, pieceNumber, pieceHash, pieceOffset, peerAddr)
				if err != nil {
					log.Println(err)
				}
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
			task.State.DownloadedBytes = task.Torrent.FileSize
			task.State.Completed = true
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

	// Return early if torrent is paused
	if task.State.Paused {
		return nil, nil
	}

	// Read piece
	pieceOffset := piece.Offset(task.Torrent.PieceSize, req.PieceNumber)
	pieceData, err := piece.Read(task.File, task.Torrent.PieceSize, pieceOffset)
	if err != nil {
		return nil, err
	}

	task.State.UploadedBytes += int64(len(pieceData))

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
		State:            &state.State{},
		AnnounceInterval: tracker.AnnounceInterval,
		AnnounceTicker:   time.NewTicker(tracker.AnnounceInterval),
		PeersRandomizer:  rand.New(rand.NewSource(time.Now().Unix())),
		WorkerPool:       worker.New(daemon.config.MaxFetchingWorkers),
		PeersNotifier:    make(chan struct{}, 1),
		PauseNotifier:    make(chan struct{}),
		ResumeNotifier:   make(chan struct{}),
		DeleteNotifier:   make(chan struct{}),
	}

	// Assign task to torrent
	daemon.mutex.Lock()
	daemon.torrents[req.Torrent.FileHash] = task
	daemon.mutex.Unlock()

	// Start torrent fetching process
	go daemon.fetching(task)
	// Keep announcing torrent to tracker
	go daemon.announcing(task)

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

	// Notify deletion, so we can exit the fetching loop
	task.DeleteNotifier <- struct{}{}

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

	// Return early if already resumed
	if !task.State.Paused {
		return nil, TorrentAlreadyResumedError
	}

	// Resume torrent
	task.State.Paused = false
	task.ResumeNotifier <- struct{}{}

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

	// Return early if already paused
	if task.State.Paused {
		return nil, TorrentAlreadyPausedError
	}

	// Pause torrent
	task.State.Paused = true
	task.PauseNotifier <- struct{}{}

	// Return to client
	return &PauseResponse{}, nil
}
