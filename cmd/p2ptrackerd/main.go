package main

import (
	"log"

	"github.com/dawidd6/p2p/pkg/version"

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
	cmdRoot.Flags().StringVarP(&config.Host, "host", "l", tracker.Host, "Tracker listening host.")
	cmdRoot.Flags().StringVarP(&config.Port, "port", "p", tracker.Port, "Tracker listening port.")
	cmdRoot.Flags().DurationVarP(&config.AnnounceInterval, "announce-interval", "n", tracker.AnnounceInterval, "Tracker announcing interval.")
	cmdRoot.Flags().DurationVarP(&config.CleanInterval, "clean-interval", "c", tracker.CleanInterval, "Tracker dangling peers cleaning interval.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
