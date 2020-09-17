package main

import (
	"log"
	"net"

	"github.com/dawidd6/p2p/version"

	"github.com/dawidd6/p2p/tracker"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	cmdRoot = &cobra.Command{
		Use:           "p2p-trackerd",
		Short:         "P2P file sharing system based on gRPC. (tracker daemon)",
		Version:       version.Version,
		RunE:          run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	listenAddr *string
)

func run(cmd *cobra.Command, args []string) error {
	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		return err
	}

	t := tracker.NewTracker()
	server := grpc.NewServer()
	service := &tracker.TrackerService{
		Announce: t.Announce,
	}

	tracker.RegisterTrackerService(server, service)
	return server.Serve(listener)
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	listenAddr = cmdRoot.Flags().StringP("address", "a", "localhost:8889", "Address on which tracker should listen.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
