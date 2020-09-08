package main

import (
	"context"
	"net"

	"github.com/dawidd6/p2p/pkg/tracker"
	"github.com/dawidd6/p2p/pkg/version"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var cmdRoot = &cobra.Command{
	Use:           "p2p-trackerd",
	Short:         "P2P tracker daemon.",
	Version:       version.Version,
	RunE:          run,
	SilenceErrors: false,
	SilenceUsage:  false,
}

func announce(ctx context.Context, in *tracker.AnnounceRequest) (*tracker.AnnounceResponse, error) {
	return &tracker.AnnounceResponse{}, nil
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	service := &tracker.TrackerService{
		Announce: announce,
	}

	tracker.RegisterTrackerService(server, service)

	return server.Serve(listener)
}

func main() {
	err := cmdRoot.Execute()
	if err != nil {
		panic(err)
	}
}
