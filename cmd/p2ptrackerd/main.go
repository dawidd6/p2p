package main

import (
	"log"

	"github.com/dawidd6/p2p/pkg/config"

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
			return tracker.Run(conf)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(0),
	}

	conf = config.Default()
)

func main() {
	cmdRoot.Flags().StringVarP(&conf.TrackerHost, "host", "l", conf.TrackerHost, "Tracker listening host.")
	cmdRoot.Flags().StringVarP(&conf.TrackerPort, "port", "p", conf.TrackerPort, "Tracker listening port.")
	cmdRoot.Flags().DurationVarP(&conf.AnnounceInterval, "announce-interval", "n", conf.AnnounceInterval, "Tracker announcing interval.")
	cmdRoot.Flags().DurationVarP(&conf.CleanInterval, "clean-interval", "c", conf.CleanInterval, "Tracker dangling peers cleaning interval.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
