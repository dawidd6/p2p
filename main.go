package main

import (
	"fmt"
	"net"

	"github.com/dawidd6/p2p/torrent"

	"github.com/dawidd6/p2p/daemon"

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
		Use:   "create FILE ...",
		Short: "Create torrent file.",
		RunE:  runCreate,
		Args:  cobra.MinimumNArgs(1),
	}
	cmdAdd = &cobra.Command{
		Use:   "add TORRENT ...",
		Short: "Add specified torrents to client.",
		RunE:  runAdd,
		Args:  cobra.MinimumNArgs(1),
	}

	daemonPort  string
	trackerPort string
	torrentName string
)

func init() {
	cmdCreate.Flags().StringVarP(&torrentName, "name", "n", "", "Torrent name.")
	cmdDaemon.Flags().StringVarP(&daemonPort, "port", "p", "0", "Port on which daemon should listen.")
	cmdTracker.Flags().StringVarP(&trackerPort, "port", "p", "0", "Port on which tracker should listen.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})
	cmdRoot.AddCommand(
		cmdAdd,
		cmdCreate,
		cmdDaemon,
		cmdTracker,
	)
}

func run(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func runCreate(cmd *cobra.Command, args []string) error {
	torr, err := torrent.CreateTorrent(torrentName, args)
	if err != nil {
		return err
	}

	return torr.Save()
}

func runAdd(cmd *cobra.Command, args []string) error {
	torrents, err := torrent.Load(args)
	if err != nil {
		return err
	}

	fmt.Println(torrents)

	return nil
}

func runDaemon(cmd *cobra.Command, args []string) error {
	listener, err := net.Listen("tcp", ":"+daemonPort)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	service := &daemon.DaemonService{}

	daemon.RegisterDaemonService(server, service)
	return server.Serve(listener)
}

func runTracker(cmd *cobra.Command, args []string) error {
	listener, err := net.Listen("tcp", ":"+trackerPort)
	if err != nil {
		return err
	}

	track := tracker.NewTracker()
	server := grpc.NewServer()
	service := &tracker.TrackerService{
		Register: track.Register,
		Lookup:   track.Lookup,
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
