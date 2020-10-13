package main

import (
	"context"
	"log"

	"github.com/dawidd6/p2p/pkg/seed"

	"github.com/dawidd6/p2p/pkg/daemon"
	"github.com/dawidd6/p2p/pkg/torrent"
	"github.com/dawidd6/p2p/pkg/tracker"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var version string

var (
	cmdRoot = &cobra.Command{
		Use:     "p2p",
		Short:   "P2P file sharing system based on gRPC.",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
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
	cmdTracker = &cobra.Command{
		Use:   "tracker",
		Short: "Run tracker.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tracker.New(*trackerAddr).Listen()
		},
		Args: cobra.ExactArgs(0),
	}
	cmdDaemon = &cobra.Command{
		Use:   "daemon",
		Short: "Run daemon.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return daemon.New(*daemonAddr).Listen()
		},
		Args: cobra.ExactArgs(0),
	}

	torrentName *string
	torrentDir  *string
	daemonAddr  *string
	seedAddr    *string
	trackerAddr *string
)

func runCreate(cmd *cobra.Command, args []string) error {
	filePaths := args

	torr, err := torrent.CreateTorrentFromFiles(*torrentName, filePaths)
	if err != nil {
		return err
	}

	return torrent.SaveTorrentToFile(torr)
}

func runAdd(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	torr, err := torrent.LoadTorrentFromFile(filePath)
	if err != nil {
		return err
	}

	err = torrent.VerifyFiles(torr, *torrentDir)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(*daemonAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	request := &daemon.AddRequest{Torrent: torr}
	client := daemon.NewDaemonClient(conn)
	_, err = client.Add(context.TODO(), request)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	trackerAddr = cmdTracker.Flags().StringP("address", "a", tracker.DefaultListenAddr, "Tracker listening address.")
	daemonAddr = cmdDaemon.Flags().StringP("address", "a", daemon.DefaultListenAddr, "Daemon listening address.")
	seedAddr = cmdDaemon.Flags().StringP("seed-address", "s", seed.DefaultListenAddr, "Seed listening address.")
	torrentName = cmdCreate.Flags().StringP("name", "n", "", "Torrent name.")
	torrentDir = cmdCreate.Flags().StringP("dir", "d", "", "Where files should be stored.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})
	cmdRoot.AddCommand(
		cmdAdd,
		cmdCreate,
		cmdTracker,
		cmdDaemon,
	)

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
