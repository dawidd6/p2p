package defaults

import "time"

const (
	TrackerListenAddress = "127.0.0.1:8889"
	DaemonListenAddress  = "127.0.0.1:8888"
	SeedListenAddress    = "0.0.0.0:44444"

	TrackerAnnounceInterval = time.Second * 30
	TrackerCleanInterval    = time.Second * 60

	TrackerConnTimeout = time.Second * 5
	DaemonConnTimeout  = time.Second * 3
	SeedConnTimeout    = time.Second * 2

	MaxWorkers = 4

	PieceSize = 256 * 1024 // 256 kB
)
