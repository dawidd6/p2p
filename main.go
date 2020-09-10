package main

import (
	"context"
	"net"

	"github.com/dawidd6/p2p/proto"

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

	daemonListenAddr  string
	trackerListenAddr string
	torrentName       string
	daemonAddr        string
)

func init() {
	cmdRoot.Flags().StringVarP(&daemonAddr, "address", "a", "localhost:8888", "Daemon address.")
	cmdCreate.Flags().StringVarP(&torrentName, "name", "n", "", "Torrent name.")
	cmdDaemon.Flags().StringVarP(&daemonListenAddr, "address", "a", "localhost:8888", "Address on which daemon should listen.")
	cmdTracker.Flags().StringVarP(&trackerListenAddr, "address", "a", "localhost:8889", "Address on which tracker should listen.")

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
	for _, filePath := range args {
		torr, err := torrent.Load(filePath)
		if err != nil {
			return err
		}

		conn, err := grpc.Dial(daemonAddr)
		if err != nil {
			return err
		}

		request := &proto.AddRequest{Torrent: torr.Torrent}
		client := proto.NewDaemonClient(conn)
		_, err = client.Add(context.TODO(), request)
		if err != nil {
			return err
		}
	}

	return nil
}

func runDaemon(cmd *cobra.Command, args []string) error {
	listener, err := net.Listen("tcp", daemonListenAddr)
	if err != nil {
		return err
	}

	d := daemon.NewDaemon()
	server := grpc.NewServer()
	service := &proto.DaemonService{
		Add: d.Add,
	}

	proto.RegisterDaemonService(server, service)
	return server.Serve(listener)
}

func runTracker(cmd *cobra.Command, args []string) error {
	listener, err := net.Listen("tcp", trackerListenAddr)
	if err != nil {
		return err
	}

	t := tracker.NewTracker()
	server := grpc.NewServer()
	service := &proto.TrackerService{
		Register: t.Register,
	}

	proto.RegisterTrackerService(server, service)
	return server.Serve(listener)
}

func main() {
	err := cmdRoot.Execute()
	if err != nil {
		panic(err)
	}
}
