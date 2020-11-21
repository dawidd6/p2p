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
		Short:   "P2P file sharing system based on gRPC (daemon).",
		Version: version.Version,
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
	cmdRoot.Flags().StringVar(&conf.DaemonHost, "host", conf.DaemonHost, "Daemon listening host.")
	cmdRoot.Flags().StringVar(&conf.DaemonPort, "port", conf.DaemonPort, "Daemon listening port.")
	cmdRoot.Flags().StringVar(&conf.SeedHost, "seed-host", conf.SeedHost, "Seed listening host.")
	cmdRoot.Flags().StringVar(&conf.SeedPort, "seed-port", conf.SeedPort, "Seed listening port.")
	cmdRoot.Flags().StringVar(&conf.DownloadsDir, "downloads-dir", conf.DownloadsDir, "Where to place downloaded files.")
	cmdRoot.Flags().StringVar(&conf.TorrentsDir, "torrents-dir", conf.TorrentsDir, "Where to place torrent files.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
