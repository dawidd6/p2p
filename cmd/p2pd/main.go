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
	cmdRoot.Flags().StringVarP(&config.Host, "host", "l", daemon.Host, "Daemon listening host.")
	cmdRoot.Flags().StringVarP(&config.Port, "port", "p", daemon.Port, "Daemon listening port.")
	cmdRoot.Flags().StringVarP(&config.SeedHost, "seed-host", "L", daemon.SeedHost, "Seed listening host.")
	cmdRoot.Flags().StringVarP(&config.SeedPort, "seed-port", "P", daemon.SeedPort, "Seed listening port.")
	cmdRoot.Flags().StringVarP(&config.DownloadsDir, "downloads-dir", "d", ".", "Where to place downloaded files.")
	cmdRoot.Flags().StringVarP(&config.SaveFilePath, "save-file", "s", daemon.SavedStateFilePath, "Where to write a state save file.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
