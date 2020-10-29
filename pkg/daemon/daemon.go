// Package daemon implements daemon with seed service
package daemon

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dawidd6/p2p/pkg/tasker"

	"github.com/dawidd6/p2p/pkg/config"

	"github.com/dawidd6/p2p/pkg/state"

	"github.com/dawidd6/p2p/pkg/hasher"
	"github.com/dawidd6/p2p/pkg/piece"

	"github.com/dawidd6/p2p/pkg/torrent"
	"github.com/dawidd6/p2p/pkg/tracker"
	"google.golang.org/grpc"
)

var (
	TorrentNotFoundError       = errors.New("torrent not found")
	TorrentAlreadyAddedError   = errors.New("torrent is already added")
	TorrentAlreadyPausedError  = errors.New("torrent is already paused")
	TorrentAlreadyResumedError = errors.New("torrent is already resumed")
	TorrentPausedError         = errors.New("torrent is paused")
	MaxSeedConnectionsError    = errors.New("max seed connections")
)

// Daemon represents daemon service (with seed)
type Daemon struct {
	conf *config.Config

	torrents      map[string]*tasker.Task
	torrentsMutex sync.RWMutex

	seedWaiter chan struct{}

	UnimplementedDaemonServer
	UnimplementedSeederServer
}

// Run starts daemon and seed servers
func Run(conf *config.Config) error {
	daemon := &Daemon{
		conf:       conf,
		torrents:   make(map[string]*tasker.Task),
		seedWaiter: make(chan struct{}, conf.MaxSeedConnections),
	}

	// Make downloads directory
	err := os.MkdirAll(conf.DownloadsDir, 0775)
	if err != nil {
		return err
	}

	// Make torrents directory
	err = os.MkdirAll(conf.TorrentsDir, 0775)
	if err != nil {
		return err
	}

	// Change to downloads directory for the whole process
	err = os.Chdir(conf.DownloadsDir)
	if err != nil {
		return err
	}

	// Restore torrents
	err = daemon.restore()
	if err != nil {
		return err
	}

	channel := make(chan error)

	// Start daemon server
	go func() {
		address := net.JoinHostPort(conf.DaemonHost, conf.DaemonPort)
		listener, err := net.Listen("tcp", address)
		if err != nil {
			channel <- err
		}

		server := grpc.NewServer()
		RegisterDaemonServer(server, daemon)
		channel <- server.Serve(listener)
	}()

	// Start seed server
	go func() {
		address := net.JoinHostPort(conf.SeedHost, conf.SeedPort)
		listener, err := net.Listen("tcp", address)
		if err != nil {
			channel <- err
		}

		server := grpc.NewServer(grpc.UnaryInterceptor(daemon.seed))
		RegisterSeederServer(server, daemon)
		channel <- server.Serve(listener)
	}()

	return <-channel
}

// restore reads saved torrent files from disk and adds them to the queue
func (daemon *Daemon) restore() error {
	glob := filepath.Join(daemon.conf.TorrentsDir, "*"+torrent.FileExtension)
	filePaths, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	if filePaths == nil {
		return nil
	}

	for i := range filePaths {
		torr, err := torrent.Load(filePaths[i])
		if err != nil {
			return err
		}

		task := tasker.New(torr, daemon.conf)

		daemon.torrentsMutex.Lock()
		daemon.torrents[torr.FileHash] = task
		daemon.torrentsMutex.Unlock()

		go daemon.fetching(task)
		go daemon.announcing(task)
	}

	return nil
}

// announce announces to tracker about having a torrent in the queue
func (daemon *Daemon) announce(task *tasker.Task) error {
	log.Println("announce", "start")
	defer log.Println("announce", "stop")

	// Connect to tracker
	conn, err := grpc.Dial(task.Torrent.TrackerAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}

	// Construct announce request
	request := &tracker.AnnounceRequest{
		FileHash: task.Torrent.FileHash,
		PeerPort: daemon.conf.SeedPort,
	}

	// Call with timeout
	ctx, cancel := context.WithTimeout(context.Background(), daemon.conf.AnnounceTimeout)
	defer cancel()

	// Announce to tracker
	client := tracker.NewTrackerClient(conn)
	response, err := client.Announce(ctx, request)
	if err != nil {
		return err
	}
	cancel()

	// Close the connection to tracker
	err = conn.Close()
	if err != nil {
		return err
	}

	// Set peer addresses
	task.PeersMutex.Lock()
	task.Peers = make(map[string]int, len(response.PeerAddresses))
	for _, peerAddr := range response.PeerAddresses {
		task.Peers[peerAddr] = 0
	}
	if len(task.PeersNotifier) == 1 {
		// Drain channel if it's full
		log.Println("announce", "draining channel")
		<-task.PeersNotifier
		log.Println("announce", "channel drained")
	}
	task.PeersMutex.Unlock()

	// Notify about new peer list
	log.Println("announce", "notifying")
	task.PeersNotifier <- struct{}{}
	log.Println("announce", "notified")

	// Set announcing interval if changed and reset the ticker
	announceInterval := time.Duration(response.AnnounceInterval) * time.Second
	if task.AnnounceInterval != announceInterval {
		log.Println("announce", "setting interval", announceInterval)
		task.AnnounceInterval = announceInterval
		task.AnnounceTicker.Reset(task.AnnounceInterval)
	}

	return nil
}

// announcing should be called in separate goroutine
func (daemon *Daemon) announcing(task *tasker.Task) {
	// Keep announcing after specified interval
	for range task.AnnounceTicker.C {
		err := daemon.announce(task)
		if err != nil {
			log.Println(err)
		}
	}
}

// fetch retrieves one piece from one peer
func (daemon *Daemon) fetch(task *tasker.Task, pieceNumber int64, pieceHash string, pieceOffset int64, peerAddr string) error {
	log.Println("fetch", pieceNumber, peerAddr)

	hash := hasher.New()

	// Try reading piece from disk
	pieceData, err := piece.Read(task.File, task.Torrent.PieceSize, pieceOffset)
	if err != nil {
		return err
	}

	// Check if read piece is correct, return if it is
	err = hash.Verify(pieceData, pieceHash)
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

	// Call with timeout
	ctx, cancel := context.WithTimeout(context.Background(), daemon.conf.FetchTimeout)
	defer cancel()

	// Get the piece from peer
	client := NewSeederClient(conn)
	response, err := client.Seed(ctx, request)
	if err != nil {
		return err
	}

	// Close the connection
	err = conn.Close()
	if err != nil {
		return err
	}

	// Got the piece wrongly, ask someone else
	err = hash.Verify(response.PieceData, pieceHash)
	if err != nil {
		return err
	}

	// Save downloaded piece on disk, can be called concurrently
	err = piece.Write(task.File, pieceOffset, response.PieceData)
	if err != nil {
		return err
	}

	// Add downloaded bytes count
	atomic.AddInt64(&task.State.DownloadedBytes, int64(len(response.PieceData)))

	return nil
}

// peer blocks until a peer is available
func (daemon *Daemon) peer(task *tasker.Task) string {
	for {
		task.PeersMutex.Lock()
		for peerAddr, peerFailures := range task.Peers {
			// If number of failures does not exceed the max, then return this peer address
			if peerFailures < daemon.conf.MaxPeerFailures {
				task.PeersMutex.Unlock()
				return peerAddr
			}
		}
		task.PeersMutex.Unlock()

		log.Println("fetch", "waiting for peers")
		// No good peer were found, wait for new list from tracker
		<-task.PeersNotifier
		log.Println("fetch", "got new peers")
	}
}

// fetching should be called in separate goroutine
func (daemon *Daemon) fetching(task *tasker.Task) {
	log.Println("add")

	// Make an initial announce to tracker
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
			peerAddr := daemon.peer(task)

			// Create worker, wait if max count already
			task.WorkerPool.Enqueue(func() {
				// Fetch one piece
				err := daemon.fetch(task, pieceNumber, pieceHash, pieceOffset, peerAddr)
				if err != nil {
					// Add peer failure
					task.PeersMutex.Lock()
					task.Peers[peerAddr]++
					task.PeersMutex.Unlock()
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

// seed is an interceptor (middleware) and it waits for available seeding slot
func (daemon *Daemon) seed(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Allow only N seed connections at the time
	select {
	case <-ctx.Done():
		return nil, MaxSeedConnectionsError
	case daemon.seedWaiter <- struct{}{}:
		res, err := handler(ctx, req)
		<-daemon.seedWaiter
		return res, err
	}
}

// Seed is called by the daemon and it shares pieces with peers
func (daemon *Daemon) Seed(ctx context.Context, req *SeedRequest) (*SeedResponse, error) {
	// Check if torrent is present in queue
	daemon.torrentsMutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.torrentsMutex.RUnlock()
	if !ok {
		return nil, TorrentNotFoundError
	}

	// Return early if torrent is paused
	if task.State.Paused {
		return nil, TorrentPausedError
	}

	// Read piece
	pieceOffset := piece.Offset(task.Torrent.PieceSize, req.PieceNumber)
	pieceData, err := piece.Read(task.File, task.Torrent.PieceSize, pieceOffset)
	if err != nil {
		return nil, err
	}

	// Add uploaded bytes count
	atomic.AddInt64(&task.State.UploadedBytes, int64(len(pieceData)))

	// Return to client
	return &SeedResponse{PieceData: pieceData}, nil
}

// Add is called by the client and it adds a new torrent to the queue
func (daemon *Daemon) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	// Check if torrent is present in queue
	daemon.torrentsMutex.RLock()
	task, ok := daemon.torrents[req.Torrent.FileHash]
	daemon.torrentsMutex.RUnlock()
	if ok {
		return nil, TorrentAlreadyAddedError
	}

	// Create new task
	task = tasker.New(req.Torrent, daemon.conf)

	// Assign task to torrent
	daemon.torrentsMutex.Lock()
	daemon.torrents[req.Torrent.FileHash] = task
	daemon.torrentsMutex.Unlock()

	// Start torrent fetching process
	go daemon.fetching(task)
	// Keep announcing torrent to tracker
	go daemon.announcing(task)

	// Save torrent file to disk
	err := torrent.Save(req.Torrent, daemon.conf.TorrentsDir, req.Torrent.FileHash)
	if err != nil {
		return nil, err
	}

	// Return to client
	return &AddResponse{}, nil
}

// Add is called by the client and it deletes a torrent from the queue
func (daemon *Daemon) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	// Check if torrent is present in queue
	daemon.torrentsMutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.torrentsMutex.RUnlock()
	if !ok {
		return nil, TorrentNotFoundError
	}

	// Delete torrent from queue
	daemon.torrentsMutex.Lock()
	delete(daemon.torrents, req.FileHash)
	daemon.torrentsMutex.Unlock()

	// Notify deletion, so we can exit the fetching loop
	if !task.State.Completed {
		task.DeleteNotifier <- struct{}{}
	}

	// Stop announcing torrent to tracker
	task.AnnounceTicker.Stop()

	// Finish any fetching workers
	task.WorkerPool.Finish()

	// Close torrent data file
	err := task.File.Close()
	if err != nil {
		return nil, err
	}

	// Remove torrent file from disk
	// TODO do it better
	filePath := filepath.Join(daemon.conf.TorrentsDir, fmt.Sprintf("%s.%s", req.FileHash, torrent.FileExtension))
	err = os.Remove(filePath)
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

// Status is called by the client and it returns data about torrent queue to the client
func (daemon *Daemon) Status(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	// Handle case when status is requested for just one torrent
	if req.FileHash != "" {
		// Check if torrent is present in queue
		daemon.torrentsMutex.RLock()
		task, ok := daemon.torrents[req.FileHash]
		daemon.torrentsMutex.RUnlock()
		if !ok {
			return nil, TorrentNotFoundError
		}

		// Return to client
		return &StatusResponse{
			States: []*state.State{task.State},
		}, nil
	}

	// Get all torrents and construct a status map
	daemon.torrentsMutex.RLock()
	index := 0
	states := make([]*state.State, len(daemon.torrents))
	for fileHash := range daemon.torrents {
		states[index] = daemon.torrents[fileHash].State
		index++
	}
	daemon.torrentsMutex.RUnlock()

	// Return to client
	return &StatusResponse{
		States: states,
	}, nil
}

// Add is called by the client and it resumes a torrent in the queue
func (daemon *Daemon) Resume(ctx context.Context, req *ResumeRequest) (*ResumeResponse, error) {
	// Check if torrent is present in queue
	daemon.torrentsMutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.torrentsMutex.RUnlock()
	if !ok {
		return nil, TorrentNotFoundError
	}

	// Return early if already resumed
	if !task.State.Paused {
		return nil, TorrentAlreadyResumedError
	}

	// Resume torrent
	task.State.Paused = false
	if !task.State.Completed {
		task.ResumeNotifier <- struct{}{}
	}

	// Return to client
	return &ResumeResponse{}, nil
}

// Add is called by the client and it pauses a torrent in the queue
func (daemon *Daemon) Pause(ctx context.Context, req *PauseRequest) (*PauseResponse, error) {
	// Check if torrent is present in queue
	daemon.torrentsMutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.torrentsMutex.RUnlock()
	if !ok {
		return nil, TorrentNotFoundError
	}

	// Return early if already paused
	if task.State.Paused {
		return nil, TorrentAlreadyPausedError
	}

	// Pause torrent
	task.State.Paused = true
	if !task.State.Completed {
		task.PauseNotifier <- struct{}{}
	}

	// Return to client
	return &PauseResponse{}, nil
}
