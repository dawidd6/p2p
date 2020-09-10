package main

import (
	"context"
	"log"
	"net"

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

	daemonListenAddr  *string
	trackerListenAddr *string
	torrentName       *string
	daemonAddr        *string
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	daemonAddr = cmdRoot.Flags().StringP("address", "a", "localhost:8888", "Daemon address.")
	torrentName = cmdCreate.Flags().StringP("name", "n", "", "Torrent name.")
	daemonListenAddr = cmdDaemon.Flags().StringP("address", "a", "localhost:8888", "Address on which daemon should listen.")
	trackerListenAddr = cmdTracker.Flags().StringP("address", "a", "localhost:8889", "Address on which tracker should listen.")

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
	filePaths := args

	torr, err := CreateTorrent(*torrentName, filePaths)
	if err != nil {
		return err
	}

	return torr.SaveTorrent()
}

func runAdd(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	torr, err := LoadTorrent(filePath)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(*daemonAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	request := &AddRequest{Torrent: torr}
	client := NewDaemonClient(conn)
	_, err = client.Add(context.TODO(), request)
	if err != nil {
		return err
	}

	return nil
}

func runDaemon(cmd *cobra.Command, args []string) error {
	listener, err := net.Listen("tcp", *daemonListenAddr)
	if err != nil {
		return err
	}

	daemon := NewDaemon()
	server := grpc.NewServer()
	service := &DaemonService{
		Add:    daemon.Add,
		Delete: daemon.Delete,
	}

	RegisterDaemonService(server, service)
	return server.Serve(listener)
}

func runTracker(cmd *cobra.Command, args []string) error {
	listener, err := net.Listen("tcp", *trackerListenAddr)
	if err != nil {
		return err
	}

	tracker := NewTracker()
	server := grpc.NewServer()
	service := &TrackerService{
		Register: tracker.Register,
	}

	RegisterTrackerService(server, service)
	return server.Serve(listener)
}

func main() {
	err := cmdRoot.Execute()
	if err != nil {
		panic(err)
	}
}
