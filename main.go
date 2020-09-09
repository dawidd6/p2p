package main

import (
	"net"

	"github.com/dawidd6/p2p/tracker"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// Set at compilation time.
var version string

var (
	cmdRoot = &cobra.Command{
		Use:           "p2p",
		Short:         "P2P file sharing system based on gRPC.",
		Version:       version,
		RunE:          run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmdDaemon = &cobra.Command{
		Use:   "daemon",
		Short: "P2P client daemon.",
		RunE:  runDaemon,
	}
	cmdTracker = &cobra.Command{
		Use:   "tracker",
		Short: "P2P tracker daemon.",
		RunE:  runTracker,
	}
	cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create metadata file.",
		RunE:  runCreate,
	}

	daemonPort  string
	trackerPort string
)

func init() {
	cmdDaemon.Flags().StringVarP(&daemonPort, "port", "p", "0", "Port on which daemon should listen.")
	cmdTracker.Flags().StringVarP(&trackerPort, "port", "p", "0", "Port on which tracker should listen.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})
	cmdRoot.AddCommand(
		cmdCreate,
		cmdDaemon,
		cmdTracker,
	)
}

func run(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func runCreate(cmd *cobra.Command, args []string) error {
	return nil
}

func runDaemon(cmd *cobra.Command, args []string) error {
	listener, err := net.Listen("tcp", ":"+daemonPort)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	service := &tracker.TrackerService{
		Announce: tracker.Announce,
	}

	tracker.RegisterTrackerService(server, service)

	return server.Serve(listener)
}

func runTracker(cmd *cobra.Command, args []string) error {
	listener, err := net.Listen("tcp", ":"+trackerPort)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	service := &tracker.TrackerService{
		Announce: tracker.Announce,
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
