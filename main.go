package main

import (
	"context"
	"log"
	"time"

	"github.com/dawidd6/p2p/pkg/defaults"

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
			config := &tracker.Config{
				Address:          *trackerAddr,
				AnnounceInterval: time.Second * 30, // TODO defaults
				CleanInterval:    time.Second * 60, // TODO defaults
			}
			return tracker.Run(config)
		},
		Args: cobra.ExactArgs(0),
	}
	cmdDaemon = &cobra.Command{
		Use:   "daemon",
		Short: "Run daemon.",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := &daemon.Config{
				DownloadsDir: *downloadsDir,
				Address:      *daemonAddr,
				SeedAddress:  *seedAddr,
			}
			return daemon.Run(config)
		},
		Args: cobra.ExactArgs(0),
	}

	trackerAddr  *string
	daemonAddr   *string
	seedAddr     *string
	pieceSize    *uint64
	downloadsDir *string
)

func runCreate(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	t, err := torrent.Create(filePath, *pieceSize)
	if err != nil {
		return err
	}

	return torrent.Save(t)
}

func runAdd(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	t, err := torrent.Load(filePath)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(*daemonAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	request := &daemon.AddRequest{Torrent: t}
	_, err = daemon.NewDaemonClient(conn).Add(context.TODO(), request)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	trackerAddr = cmdTracker.Flags().StringP("address", "a", defaults.TrackerListenAddress, "Tracker listening address.")
	daemonAddr = cmdDaemon.Flags().StringP("address", "a", defaults.DaemonListenAddress, "Daemon listening address.")
	seedAddr = cmdDaemon.Flags().StringP("seed-address", "s", defaults.SeedListenAddress, "Seed listening address.")
	downloadsDir = cmdDaemon.Flags().StringP("downloads-dir", "d", ".", "Where to place downloaded files.")
	pieceSize = cmdCreate.Flags().Uint64P("piece-size", "s", defaults.PieceSize, "Piece size.")

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
