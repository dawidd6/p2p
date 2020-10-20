package main

import (
	"log"

	"github.com/dawidd6/p2p/pkg/version"

	"github.com/dawidd6/p2p/pkg/defaults"

	"github.com/dawidd6/p2p/pkg/tracker"

	"github.com/spf13/cobra"
)

var (
	cmdRoot = &cobra.Command{
		Use:     "p2ptrackerd",
		Short:   "P2P file sharing system based on gRPC.",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return tracker.Run(config)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(0),
	}

	config = &tracker.Config{}
)

func main() {
	cmdRoot.Flags().StringVarP(&config.ListenAddress, "listen-address", "l", defaults.TrackerListenAddress, "Tracker listening address.")
	cmdRoot.Flags().DurationVarP(&config.AnnounceInterval, "announce-interval", "n", defaults.TrackerAnnounceInterval, "Tracker announcing interval.")
	cmdRoot.Flags().DurationVarP(&config.CleanInterval, "clean-interval", "c", defaults.TrackerCleanInterval, "Tracker dangling peers cleaning interval.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
