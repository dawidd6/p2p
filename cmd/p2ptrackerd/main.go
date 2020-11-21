package main

import (
	"log"

	"github.com/dawidd6/p2p/pkg/config"

	"github.com/dawidd6/p2p/cmd"

	"github.com/dawidd6/p2p/pkg/tracker"

	"github.com/spf13/cobra"
)

var (
	cmdRoot = &cobra.Command{
		Use:     "p2ptrackerd",
		Short:   "P2P file sharing system based on gRPC (tracker).",
		Version: cmd.Version,
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
	cmdRoot.Flags().StringVar(&conf.TrackerHost, "host", conf.TrackerHost, "Tracker listening host.")
	cmdRoot.Flags().StringVar(&conf.TrackerPort, "port", conf.TrackerPort, "Tracker listening port.")
	cmdRoot.Flags().StringVar(&conf.DBHost, "db-host", conf.DBHost, "Database listening host.")
	cmdRoot.Flags().StringVar(&conf.DBPort, "db-port", conf.DBPort, "Database listening port.")
	cmdRoot.Flags().DurationVar(&conf.AnnounceInterval, "announce-interval", conf.AnnounceInterval, "Tracker announcing interval.")
	cmdRoot.Flags().DurationVar(&conf.CleanInterval, "clean-interval", conf.CleanInterval, "Tracker dangling peers cleaning interval.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
