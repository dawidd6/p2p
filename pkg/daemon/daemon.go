// Package daemon implements daemon with seed service
package daemon

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	conf          *config.Config
	torrents      map[string]*tasker.Task
	torrentsMutex sync.RWMutex
	seedWaiter    chan struct{}

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

		log.Println("Listening on", address, "(daemon)")

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

		log.Println("Listening on", address, "(seed)")

		server := grpc.NewServer(grpc.UnaryInterceptor(daemon.seed))
		RegisterSeederServer(server, daemon)
		channel <- server.Serve(listener)
	}()

	return <-channel
}

// restore reads saved torrent files from disk and adds them back to the queue
func (daemon *Daemon) restore() error {
	// Find all torrent files on disk in directory
	torrentGlob := filepath.Join(daemon.conf.TorrentsDir, "*"+torrent.FileExtension)
	torrentFilePaths, err := filepath.Glob(torrentGlob)
	if err != nil {
		return err
	}

	// Just exit if no files were found
	if torrentFilePaths == nil {
		return nil
	}

	log.Println("Found", len(torrentFilePaths), "torrents to restore")

	// Add every found torrent back to the queue
	for i := range torrentFilePaths {
		// Open torrent file
		torrentFile, err := os.OpenFile(torrentFilePaths[i], os.O_RDONLY, 0666)
		if err != nil {
			return err
		}

		// Read the torrent from file
		torr, err := torrent.Read(torrentFile)
		if err != nil {
			return err
		}

		// Open state file
		stateFilePath := strings.TrimSuffix(torrentFilePaths[i], torrent.FileExtension) + state.FileExtension
		stateFile, err := os.OpenFile(stateFilePath, os.O_RDONLY, 0666)
		if err != nil {
			return err
		}

		// Read the state from file
		stat, err := state.Read(stateFile)
		if err != nil {
			return err
		}

		// Close torrent file, it is later reopened
		err = torrentFile.Close()
		if err != nil {
			return err
		}

		// Close state file, it is later reopened
		err = stateFile.Close()
		if err != nil {
			return err
		}

		// Add torrent to the queue
		err = daemon.add(torr, stat)
		if err != nil {
			return err
		}
	}

	return nil
}

// saving should be called in separate goroutine
func (daemon *Daemon) saving(task *tasker.Task) {
	for {
		err := task.StateFile.Truncate(0)
		if err != nil {
			log.Println(err)
		}

		_, err = task.StateFile.Seek(0, 0)
		if err != nil {
			log.Println(err)
		}

		err = state.Write(task.StateFile, task.State)
		if err != nil {
			log.Println(err)
		}

		select {
		// Exit if stop was requested
		case <-task.SaveNotifier:
			log.Println("Stopping saving state of torrent", task.Torrent.FileHash)
			return
		// Keep announcing after specified interval
		case <-task.SaveTicker.C:
		}
	}
}

// announce announces to tracker about having a torrent in the queue
func (daemon *Daemon) announce(task *tasker.Task) error {
	log.Println("Announcing", task.Torrent.FileHash, "to", task.Torrent.TrackerAddress)

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

	// Close the connection to tracker
	err = conn.Close()
	if err != nil {
		return err
	}

	// Set peer addresses and failures
	task.PeersMutex.Lock()
	task.Peers = make(map[string]int, len(response.PeerAddresses))
	for _, peerAddr := range response.PeerAddresses {
		task.Peers[peerAddr] = 0
	}
	task.PeersMutex.Unlock()

	// Notify about new peer list
	task.PeersAvailable.Broadcast()

	// Set announcing interval if changed and reset the ticker
	announceInterval := time.Duration(response.AnnounceInterval) * time.Second
	if task.AnnounceInterval != announceInterval {
		log.Println("Setting announce interval", announceInterval, "for", task.Torrent.FileHash)
		task.AnnounceInterval = announceInterval
		task.AnnounceTicker.Reset(task.AnnounceInterval)
	}

	// Set state peer count
	task.StateMutex.Lock()
	task.State.PeersCount = int64(len(task.Peers))
	task.StateMutex.Unlock()

	return nil
}

// announcing should be called in separate goroutine
func (daemon *Daemon) announcing(task *tasker.Task) {
	for {
		err := daemon.announce(task)
		if err != nil {
			log.Println(err)
			// Reset peer count
			task.StateMutex.Lock()
			task.State.PeersCount = 0
			task.StateMutex.Unlock()
		}

		select {
		// Exit if stop was requested
		case <-task.AnnounceNotifier:
			log.Println("Stopping announcing of torrent", task.Torrent.FileHash)
			return
		// Keep announcing after specified interval
		case <-task.AnnounceTicker.C:
		}
	}
}

// fetch retrieves one piece from one peer
func (daemon *Daemon) fetch(task *tasker.Task, pieceNumber int64, pieceHash string, pieceOffset int64, peerAddr string) error {
	log.Println("Fetching", pieceNumber, "piece of", task.Torrent.FileHash, "from", peerAddr)

	hash := hasher.New()

	// Try reading piece from disk
	pieceData, err := piece.Read(task.DataFile, task.Torrent.PieceSize, pieceOffset)
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
	err = piece.Write(task.DataFile, pieceOffset, response.PieceData)
	if err != nil {
		return err
	}

	// Add downloaded bytes count and compute ratio
	task.StateMutex.Lock()
	task.State.DownloadedBytes += int64(len(response.PieceData))
	task.State.Ratio = float32(task.State.UploadedBytes) / float32(task.State.DownloadedBytes)
	task.StateMutex.Unlock()

	return nil
}

// peer blocks until a valid peer is available
func (daemon *Daemon) peer(task *tasker.Task) string {
	for {
		task.PeersMutex.RLock()
		for peerAddr, peerFailures := range task.Peers {
			// If number of failures does not exceed the max, then return this peer address
			if peerFailures < daemon.conf.MaxPeerFailures {
				task.PeersMutex.RUnlock()
				return peerAddr
			}
		}
		task.PeersMutex.RUnlock()

		// No good peer were found, wait for new list from tracker
		log.Println("Waiting for peers for", task.Torrent.FileHash)
		task.PeersAvailable.L.Lock()
		task.PeersAvailable.Wait()
		task.PeersAvailable.L.Unlock()
	}
}

// fetching should be called in separate goroutine
func (daemon *Daemon) fetching(task *tasker.Task) {
	// Check if torrent is already downloaded fully
	err := torrent.Verify(task.DataFile, task.Torrent)
	if err == nil {
		log.Println("Already completed", task.Torrent.FileHash)
		task.StateMutex.Lock()
		task.State.DownloadedBytes = task.Torrent.FileSize
		task.State.Completed = true
		task.StateMutex.Unlock()
		return
	} else {
		task.StateMutex.Lock()
		task.State.Completed = false
		task.StateMutex.Unlock()
	}

	// If initial torrent state is paused, wait for resume
	if task.State.Paused {
		log.Println("Started", task.Torrent.FileHash, "in paused state")
		<-task.ResumeNotifier
	}

	// Start a fetch loop
	for {
		// Start a worker pool
		log.Println("Starting worker pool")
		task.WorkerPool.Start()

		// Loop over torrent piece hashes
		for i := range task.Torrent.PieceHashes {
			pieceNumber := int64(i)
			pieceHash := task.Torrent.PieceHashes[pieceNumber]
			pieceOffset := piece.Offset(task.Torrent.PieceSize, pieceNumber)

			// Pause execution if desired or exit from function if torrent is deleted
			select {
			case <-task.PauseNotifier:
				log.Println("Paused", task.Torrent.FileHash)
				<-task.ResumeNotifier
				log.Println("Resumed", task.Torrent.FileHash)
			case <-task.DeleteNotifier:
				log.Println("Deleted", task.Torrent.FileHash)
				return
			default:
			}

			// Send a job to a free worker, wait if workers are busy
			task.WorkerPool.Enqueue(func() {
				// Get random peer address, wait if no peers available
				peerAddr := daemon.peer(task)

				// Fetch one piece
				err := daemon.fetch(task, pieceNumber, pieceHash, pieceOffset, peerAddr)
				if err != nil {
					// Add peer failure
					log.Println("Adding peer", peerAddr, "failure for", task.Torrent.FileHash)
					task.PeersMutex.Lock()
					task.Peers[peerAddr]++
					task.PeersMutex.Unlock()
				}
			})
		}

		// Wait for all fetch workers to complete
		log.Println("Stopping worker pool for", task.Torrent.FileHash)
		task.WorkerPool.Stop()

		// Verify if all downloaded pieces are correct
		err := torrent.Verify(task.DataFile, task.Torrent)
		if err != nil {
			log.Println("Torrent", task.Torrent.FileHash, "has not been downloaded properly")
			log.Println("Retrying torrent", task.Torrent.FileHash, "download")
			task.StateMutex.Lock()
			task.State.Completed = false
			task.StateMutex.Unlock()
			continue
		} else {
			log.Println("Torrent", task.Torrent.FileHash, "download complete")
			task.StateMutex.Lock()
			task.State.DownloadedBytes = task.Torrent.FileSize
			task.State.Completed = true
			task.StateMutex.Unlock()
			return
		}
	}
}

// seed is an interceptor (middleware) and it waits for available seeding slot
func (daemon *Daemon) seed(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	select {
	// Exit if timeout
	case <-ctx.Done():
		return nil, MaxSeedConnectionsError
	// Allow only N seed connections at the time
	case daemon.seedWaiter <- struct{}{}:
		res, err := handler(ctx, req)
		<-daemon.seedWaiter
		return res, err
	}
}

// add adds a torrent to the queue
func (daemon *Daemon) add(torr *torrent.Torrent, stat *state.State) error {
	var err error

	// Check if torrent is present in queue
	daemon.torrentsMutex.RLock()
	task, ok := daemon.torrents[torr.FileHash]
	daemon.torrentsMutex.RUnlock()
	if ok {
		return TorrentAlreadyAddedError
	}

	// Create new task
	task = tasker.New(torr, stat, daemon.conf)

	log.Println("Adding new torrent", task.Torrent.FileHash)

	// Open or create the torrent data file
	dataFilePath := task.Torrent.FileName
	task.DataFile, err = os.OpenFile(dataFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	// Open the torrent file
	torrentFilePath := filepath.Join(daemon.conf.TorrentsDir, task.Torrent.FileHash+torrent.FileExtension)
	task.TorrentFile, err = os.OpenFile(torrentFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	// Open the state file
	stateFilePath := filepath.Join(daemon.conf.TorrentsDir, task.Torrent.FileHash+state.FileExtension)
	task.StateFile, err = os.OpenFile(stateFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	// Save torrent file to disk
	err = torrent.Write(task.TorrentFile, task.Torrent)
	if err != nil {
		return err
	}

	// Save torrent state to disk
	err = state.Write(task.StateFile, task.State)
	if err != nil {
		return err
	}

	// Assign task to torrent
	daemon.torrentsMutex.Lock()
	daemon.torrents[task.Torrent.FileHash] = task
	daemon.torrentsMutex.Unlock()

	// Keep announcing torrent to tracker
	go daemon.announcing(task)
	// Keep saving torrent state to file
	go daemon.saving(task)
	// Start torrent fetching process
	go daemon.fetching(task)

	return nil
}

func (daemon *Daemon) delete(fileHash string, withData bool) error {
	var err error

	// Check if torrent is present in queue
	daemon.torrentsMutex.RLock()
	task, ok := daemon.torrents[fileHash]
	daemon.torrentsMutex.RUnlock()
	if !ok {
		return TorrentNotFoundError
	}

	// Delete torrent from queue
	daemon.torrentsMutex.Lock()
	delete(daemon.torrents, task.Torrent.FileHash)
	daemon.torrentsMutex.Unlock()

	log.Println("Deleting torrent", task.Torrent.FileHash, "with data:", withData)

	// Notify deletion, so we can exit the fetching loop
	task.StateMutex.RLock()
	if !task.State.Completed {
		task.DeleteNotifier <- struct{}{}
	}
	task.StateMutex.RUnlock()

	// Stop announcing torrent to tracker
	task.AnnounceNotifier <- struct{}{}
	task.AnnounceTicker.Stop()

	// Stop saving torrent state to file
	task.SaveNotifier <- struct{}{}
	task.SaveTicker.Stop()

	// Finish any fetching workers
	task.WorkerPool.Stop()

	// Remove downloaded torrent data if desired
	if withData {
		err = os.Remove(task.DataFile.Name())
		if err != nil {
			return err
		}
	}

	// Close torrent data file
	err = task.DataFile.Close()
	if err != nil {
		return err
	}

	// Remove torrent file
	err = os.Remove(task.TorrentFile.Name())
	if err != nil {
		return err
	}

	// Close torrent file
	err = task.TorrentFile.Close()
	if err != nil {
		return err
	}

	// Remove state file
	err = os.Remove(task.StateFile.Name())
	if err != nil {
		return err
	}

	// Close state file
	err = task.StateFile.Close()
	if err != nil {
		return err
	}

	return nil
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

	log.Println("Sending piece", req.PieceNumber, "of", req.FileHash)

	// Read piece
	pieceOffset := piece.Offset(task.Torrent.PieceSize, req.PieceNumber)
	pieceData, err := piece.Read(task.DataFile, task.Torrent.PieceSize, pieceOffset)
	if err != nil {
		return nil, err
	}

	// Add uploaded bytes count and compute ratio
	task.StateMutex.Lock()
	task.State.UploadedBytes += int64(len(pieceData))
	task.State.Ratio = float32(task.State.UploadedBytes) / float32(task.State.DownloadedBytes)
	task.StateMutex.Unlock()

	// Return to client
	return &SeedResponse{PieceData: pieceData}, nil
}

// Add is called by the client and it adds a new torrent to the queue
func (daemon *Daemon) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	// Add torrent
	err := daemon.add(req.Torrent, nil)
	if err != nil {
		return nil, err
	}

	// Return to client
	return &AddResponse{}, nil
}

// Delete is called by the client and it deletes a torrent from the queue
func (daemon *Daemon) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	// Delete torrent
	err := daemon.delete(req.FileHash, req.WithData)
	if err != nil {
		return nil, err
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

		// Lock, return state, unlock
		task.StateMutex.RLock()
		defer task.StateMutex.RUnlock()

		// Return to client
		return &StatusResponse{
			States: []*state.State{task.State},
		}, nil
	}

	// Get all torrents and construct a status map
	daemon.torrentsMutex.RLock()
	index := 0
	states := make([]*state.State, len(daemon.torrents))
	for _, task := range daemon.torrents {
		task.StateMutex.RLock()
		states[index] = task.State
		task.StateMutex.RUnlock()
		index++
	}
	daemon.torrentsMutex.RUnlock()

	// Return to client
	return &StatusResponse{
		States: states,
	}, nil
}

// Resume is called by the client and it resumes a torrent in the queue
func (daemon *Daemon) Resume(ctx context.Context, req *ResumeRequest) (*ResumeResponse, error) {
	// Check if torrent is present in queue
	daemon.torrentsMutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.torrentsMutex.RUnlock()
	if !ok {
		return nil, TorrentNotFoundError
	}

	// Lock state for further operations on it
	task.StateMutex.Lock()
	defer task.StateMutex.Unlock()

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

// Pause is called by the client and it pauses a torrent in the queue
func (daemon *Daemon) Pause(ctx context.Context, req *PauseRequest) (*PauseResponse, error) {
	// Check if torrent is present in queue
	daemon.torrentsMutex.RLock()
	task, ok := daemon.torrents[req.FileHash]
	daemon.torrentsMutex.RUnlock()
	if !ok {
		return nil, TorrentNotFoundError
	}

	// Lock state for further operations on it
	task.StateMutex.Lock()
	defer task.StateMutex.Unlock()

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
