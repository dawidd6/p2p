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

	Torrent *torrent.Torrent
	State   *state.State

	ResumeNotifier chan struct{}
	PauseNotifier  chan struct{}
	DeleteNotifier chan struct{}

	PeersNotifier chan struct{}
	PeersMutex    sync.Mutex
	Peers         map[string]int

	AnnounceTicker   *time.Ticker
	AnnounceInterval time.Duration

	WorkerPool *worker.Pool
}

func New(torr *torrent.Torrent, conf *config.Config) *Task {
	return &Task{
		Torrent: torr,
		State: &state.State{
			FileHash: torr.FileHash,
			FileName: torr.FileName,
		},

		ResumeNotifier: make(chan struct{}),
		PauseNotifier:  make(chan struct{}),
		DeleteNotifier: make(chan struct{}),

		PeersNotifier: make(chan struct{}, 1),

		AnnounceTicker:   time.NewTicker(conf.AnnounceInterval),
		AnnounceInterval: conf.AnnounceInterval,

		WorkerPool: worker.New(conf.MaxFetchConnections),
	}
}
