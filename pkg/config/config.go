// Package config includes configuration options for the whole system
package config

import (
	"os"
	"path/filepath"
	"time"
)

// Config struct holds all configurable aspects of the system
type Config struct {
	TrackerHost string
	TrackerPort string

	AnnounceTolerance time.Duration
	AnnounceInterval  time.Duration
	CleanInterval     time.Duration

	DaemonHost string
	DaemonPort string

	SeedHost string
	SeedPort string

	DownloadsDir string
	TorrentsDir  string

	MaxPeerFailures     int
	MaxSeedConnections  int
	MaxFetchConnections int

	PieceSize int64

	AnnounceTimeout time.Duration
	FetchTimeout    time.Duration
	CallTimeout     time.Duration

	DeleteWithData bool
}

// Default returns an instance of config with default values set
func Default() *Config {
	return &Config{
		TrackerHost: "127.0.0.1",
		TrackerPort: "8889",

		AnnounceTolerance: time.Second * 10,
		AnnounceInterval:  time.Second * 30,
		CleanInterval:     time.Second * 60,

		DaemonHost: "127.0.0.1",
		DaemonPort: "8888",

		SeedHost: "0.0.0.0",
		SeedPort: "44444",

		DownloadsDir: ".",
		TorrentsDir:  filepath.Join(os.Getenv("HOME"), ".p2p"),

		MaxPeerFailures:     4,
		MaxSeedConnections:  4,
		MaxFetchConnections: 4,

		PieceSize: 256 * 1024,

		AnnounceTimeout: time.Second * 5,
		FetchTimeout:    time.Second * 3,
		CallTimeout:     time.Second * 5,

		DeleteWithData: false,
	}
}
