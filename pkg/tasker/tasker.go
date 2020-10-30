// Package tasker includes a Task struct
package tasker

import (
	"os"
	"sync"
	"time"

	"github.com/dawidd6/p2p/pkg/notifier"

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

	ResumeNotifier *notifier.Notifier
	PauseNotifier  *notifier.Notifier
	DeleteNotifier *notifier.Notifier

	PeersNotifier *notifier.Notifier
	PeersMutex    sync.Mutex
	Peers         map[string]int

	AnnounceNotifier *notifier.Notifier
	AnnounceTicker   *time.Ticker
	AnnounceInterval time.Duration

	WorkerPool *worker.Pool
}

// New returns a new Task instance
func New(torr *torrent.Torrent, conf *config.Config) *Task {
	return &Task{
		Torrent: torr,
		State: &state.State{
			FileHash: torr.FileHash,
			FileName: torr.FileName,
		},

		ResumeNotifier: notifier.NewBlocking(),
		PauseNotifier:  notifier.NewBlocking(),
		DeleteNotifier: notifier.NewBlocking(),

		PeersNotifier: notifier.NewNotBlocking(),

		AnnounceNotifier: notifier.NewNotBlocking(),
		AnnounceTicker:   time.NewTicker(conf.AnnounceInterval),
		AnnounceInterval: conf.AnnounceInterval,

		WorkerPool: worker.New(conf.MaxFetchConnections),
	}
}
