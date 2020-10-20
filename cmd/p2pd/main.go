package main

import (
	"log"

	"github.com/dawidd6/p2p/pkg/defaults"
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
	cmdRoot.Flags().StringVarP(&config.ListenAddress, "listen-address", "l", defaults.DaemonListenAddress, "Daemon listening address.")
	cmdRoot.Flags().StringVarP(&config.SeedListenAddress, "seed-listen-address", "s", defaults.SeedListenAddress, "Seed listening address.")
	cmdRoot.Flags().StringVarP(&config.DownloadsDir, "downloads-dir", "d", ".", "Where to place downloaded files.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
