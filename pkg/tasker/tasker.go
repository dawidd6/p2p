// Package tasker includes a Task struct
package tasker

import (
	"os"
	"sync"
	"time"

	"github.com/dawidd6/p2p/pkg/config"

	"github.com/dawidd6/p2p/pkg/state"
	"github.com/dawidd6/p2p/pkg/torrent"
	"github.com/dawidd6/p2p/pkg/worker"
)

// Task holds all the info needed for a torrent process
type Task struct {
	TorrentFile *os.File
	DataFile    *os.File
	StateFile   *os.File

	Torrent *torrent.Torrent

	State      *state.State
	StateMutex sync.RWMutex

	ResumeNotifier chan struct{}
	PauseNotifier  chan struct{}
	DeleteNotifier chan struct{}

	PeersAvailable *sync.Cond
	PeersMutex     sync.RWMutex
	Peers          map[string]int

	AnnounceNotifier chan struct{}
	AnnounceTicker   *time.Ticker
	AnnounceInterval time.Duration

	SaveNotifier chan struct{}
	SaveTicker   *time.Ticker
	SaveInterval time.Duration

	WorkerPool *worker.Pool
}

// New returns a new Task instance
func New(torr *torrent.Torrent, stat *state.State, conf *config.Config) *Task {
	task := &Task{
		Torrent: torr,

		State: &state.State{
			FileHash: torr.FileHash,
			FileName: torr.FileName,
		},

		ResumeNotifier: make(chan struct{}, 1),
		PauseNotifier:  make(chan struct{}, 1),
		DeleteNotifier: make(chan struct{}, 1),

		PeersAvailable: sync.NewCond(&sync.Mutex{}),

		AnnounceNotifier: make(chan struct{}, 1),
		AnnounceTicker:   time.NewTicker(conf.AnnounceInterval),
		AnnounceInterval: conf.AnnounceInterval,

		SaveNotifier: make(chan struct{}, 1),
		SaveTicker:   time.NewTicker(conf.SaveInterval),
		SaveInterval: conf.AnnounceInterval,

		WorkerPool: worker.New(conf.MaxFetchConnections),
	}

	// Used for state restoring
	if stat != nil {
		task.State = stat
	}

	return task
}
