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

	DBHost string
	DBPort string

	DBCheckInterval   time.Duration
	AnnounceTolerance time.Duration
	AnnounceInterval  time.Duration

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

	SaveInterval time.Duration

	DeleteWithData bool
	PrintJSON      bool
}

// Default returns an instance of config with default values set
func Default() *Config {
	return &Config{
		TrackerHost: "127.0.0.1",
		TrackerPort: "8889",

		DBHost: "127.0.0.1",
		DBPort: "6379",

		DBCheckInterval:   time.Second * 10,
		AnnounceTolerance: time.Second * 10,
		AnnounceInterval:  time.Second * 30,

		DaemonHost: "127.0.0.1",
		DaemonPort: "8888",

		SeedHost: "0.0.0.0",
		SeedPort: "44444",

		DownloadsDir: ".",
		TorrentsDir:  filepath.Join(os.Getenv("HOME"), ".p2p"),

		MaxPeerFailures:     2, // In reality it is 5. Best to keep this as is, it's a bit racy
		MaxSeedConnections:  4,
		MaxFetchConnections: 4,

		PieceSize: 256 * 1024,

		AnnounceTimeout: time.Second * 5,
		FetchTimeout:    time.Second * 3,
		CallTimeout:     time.Second * 5,

		SaveInterval: time.Second * 3,

		DeleteWithData: false,
		PrintJSON:      false,
	}
}
