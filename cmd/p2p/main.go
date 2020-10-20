package main

import (
	"context"
	"log"

	"github.com/dawidd6/p2p/pkg/version"

	"github.com/dawidd6/p2p/pkg/defaults"

	"github.com/dawidd6/p2p/pkg/daemon"
	"github.com/dawidd6/p2p/pkg/torrent"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	cmdRoot = &cobra.Command{
		Use:     "p2p",
		Short:   "P2P file sharing system based on gRPC.",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmdCreate = &cobra.Command{
		Use:   "create FILE",
		Short: "Create torrent file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			t, err := torrent.Create(filePath, *pieceSize, *trackerAddr)
			if err != nil {
				return err
			}

			return torrent.Save(t)
		},
		Args: cobra.ExactArgs(1),
	}

	cmdAdd = &cobra.Command{
		Use:   "add TORRENT",
		Short: "Add specified torrents to client.",
		RunE: func(cmd *cobra.Command, args []string) error {
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
		},
		Args: cobra.ExactArgs(1),
	}

	trackerAddr *string
	daemonAddr  *string
	pieceSize   *uint64
)

func main() {
	trackerAddr = cmdCreate.Flags().StringP("tracker-address", "t", defaults.TrackerListenAddress, "Tracker address.")
	pieceSize = cmdCreate.Flags().Uint64P("piece-size", "s", defaults.PieceSize, "Piece size.")
	daemonAddr = cmdAdd.Flags().StringP("listen-address", "l", defaults.DaemonListenAddress, "Daemon listening address.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})
	cmdRoot.AddCommand(
		cmdAdd,
		cmdCreate,
	)

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
