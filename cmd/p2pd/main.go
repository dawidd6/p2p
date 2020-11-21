package main

import (
	"log"

	"github.com/dawidd6/p2p/pkg/config"

	"github.com/dawidd6/p2p/cmd"

	"github.com/dawidd6/p2p/pkg/daemon"
	"github.com/spf13/cobra"
)

var (
	cmdRoot = &cobra.Command{
		Use:     "p2pd",
		Short:   "P2P file sharing system based on gRPC (daemon)",
		Version: cmd.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return daemon.Run(conf)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(0),
	}

	conf = config.Default()
)

func main() {
	cmdRoot.Flags().StringVar(&conf.DaemonHost, "host", conf.DaemonHost, "daemon listening host")
	cmdRoot.Flags().StringVar(&conf.DaemonPort, "port", conf.DaemonPort, "daemon listening port")
	cmdRoot.Flags().StringVar(&conf.SeedHost, "seed-host", conf.SeedHost, "seed listening host")
	cmdRoot.Flags().StringVar(&conf.SeedPort, "seed-port", conf.SeedPort, "seed listening port")
	cmdRoot.Flags().StringVar(&conf.DownloadsDir, "downloads-dir", conf.DownloadsDir, "where to place downloaded files")
	cmdRoot.Flags().StringVar(&conf.TorrentsDir, "torrents-dir", conf.TorrentsDir, "where to place torrent files")
	cmdRoot.Flags().DurationVar(&conf.SaveInterval, "save-interval", conf.SaveInterval, "how often to save torrent state")
	cmdRoot.Flags().IntVar(&conf.MaxPeerFailures, "max-peer-failures", conf.MaxPeerFailures, "max peer failures")
	cmdRoot.Flags().IntVar(&conf.MaxSeedConnections, "max-seed-connections", conf.MaxSeedConnections, "max concurrent seeding connections per torrent")
	cmdRoot.Flags().IntVar(&conf.MaxFetchConnections, "max-fetch-connections", conf.MaxFetchConnections, "max concurrent fetching connections per torrent")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
