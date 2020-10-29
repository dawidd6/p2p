package main

import (
	"log"

	"github.com/dawidd6/p2p/pkg/config"

	"github.com/dawidd6/p2p/pkg/version"

	"github.com/dawidd6/p2p/pkg/daemon"
	"github.com/spf13/cobra"
)

var (
	cmdRoot = &cobra.Command{
		Use:     "p2pd",
		Short:   "P2P file sharing system based on gRPC.",
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
	cmdRoot.Flags().StringVarP(&conf.DaemonHost, "host", "l", conf.DaemonHost, "Daemon listening host.")
	cmdRoot.Flags().StringVarP(&conf.DaemonPort, "port", "p", conf.DaemonPort, "Daemon listening port.")
	cmdRoot.Flags().StringVarP(&conf.SeedHost, "seed-host", "L", conf.SeedHost, "Seed listening host.")
	cmdRoot.Flags().StringVarP(&conf.SeedPort, "seed-port", "P", conf.SeedPort, "Seed listening port.")
	cmdRoot.Flags().StringVarP(&conf.DownloadsDir, "downloads-dir", "d", conf.DownloadsDir, "Where to place downloaded files.")
	cmdRoot.Flags().StringVarP(&conf.TorrentsDir, "torrents-dir", "t", conf.TorrentsDir, "Where to place torrent files.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
