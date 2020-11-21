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
		Short:   "P2P file sharing system based on gRPC (tracker)",
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
	cmdRoot.Flags().StringVar(&conf.TrackerHost, "host", conf.TrackerHost, "tracker listening host")
	cmdRoot.Flags().StringVar(&conf.TrackerPort, "port", conf.TrackerPort, "tracker listening port")
	cmdRoot.Flags().StringVar(&conf.DBHost, "db-host", conf.DBHost, "database listening host")
	cmdRoot.Flags().StringVar(&conf.DBPort, "db-port", conf.DBPort, "database listening port")
	cmdRoot.Flags().DurationVar(&conf.AnnounceInterval, "announce-interval", conf.AnnounceInterval, "tracker announcing interval")
	cmdRoot.Flags().DurationVar(&conf.AnnounceTolerance, "announce-tolerance", conf.AnnounceTolerance, "tracker announcing tolerance")
	cmdRoot.Flags().DurationVar(&conf.DBCheckInterval, "db-check-interval", conf.DBCheckInterval, "database checking interval")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
