package main

import (
	"log"

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
			return daemon.Run(config)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(0),
	}

	config = &daemon.Config{}
)

func main() {
	cmdRoot.Flags().StringVarP(&config.ListenAddress, "listen-address", "l", daemon.ListenAddress, "Daemon listening address.")
	cmdRoot.Flags().StringVarP(&config.SeedListenAddress, "seed-listen-address", "s", daemon.SeedListenAddress, "Seed listening address.")
	cmdRoot.Flags().StringVarP(&config.DownloadsDir, "downloads-dir", "d", daemon.DownloadsDir, "Where to place downloaded files.")
	cmdRoot.Flags().IntVarP(&config.MaxFetchingWorkers, "max-fetch-workers", "w", daemon.MaxFetchingWorkers, "Max number of fetch workers.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
