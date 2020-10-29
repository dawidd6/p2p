package config

import (
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	TrackerHost string
	TrackerPort string

	AnnounceInterval time.Duration
	CleanInterval    time.Duration

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
}

// Default returns an instance of config with default values set
func Default() *Config {
	return &Config{
		TrackerHost: "0.0.0.0",
		TrackerPort: "8889",

		AnnounceInterval: time.Second * 30,
		CleanInterval:    time.Second * 60,

		DaemonHost: "127.0.0.1",
		DaemonPort: "8888",

		SeedHost: "0.0.0.0",
		SeedPort: "44444",

		DownloadsDir: ".",
		TorrentsDir:  filepath.Join(os.Getenv("HOME"), ".p2p"),

		MaxPeerFailures:     4,
		MaxSeedConnections:  4,
		MaxFetchConnections: 4,

		PieceSize: 256 * 1000,

		AnnounceTimeout: time.Second * 5,
		FetchTimeout:    time.Second * 3,
	}
}
