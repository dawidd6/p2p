package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dawidd6/p2p/pkg/piece"
	"github.com/dawidd6/p2p/pkg/tracker"

	"github.com/dawidd6/p2p/pkg/version"

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

			file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
			if err != nil {
				return err
			}

			t, err := torrent.Create(file, *pieceSize, *trackerAddr)
			if err != nil {
				return err
			}

			err = torrent.Verify(t, file)
			if err != nil {
				return err
			}

			err = torrent.Save(t)
			if err != nil {
				return err
			}

			return file.Close()
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

			request := &daemon.AddRequest{
				Torrent: t,
			}

			client := daemon.NewDaemonClient(conn)
			_, err = client.Add(context.Background(), request)
			if err != nil {
				return err
			}

			return conn.Close()
		},
		Args: cobra.ExactArgs(1),
	}

	cmdDelete = &cobra.Command{
		Use:   "delete FILE_HASH",
		Short: "Delete specified torrent from client by file hash.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileHash := args[0]

			conn, err := grpc.Dial(*daemonAddr, grpc.WithInsecure())
			if err != nil {
				return err
			}

			request := &daemon.DeleteRequest{
				FileHash: fileHash,
				WithData: *withData,
			}

			client := daemon.NewDaemonClient(conn)
			_, err = client.Delete(context.Background(), request)
			if err != nil {
				return err
			}

			return conn.Close()
		},
		Args: cobra.ExactArgs(1),
	}

	cmdStatus = &cobra.Command{
		Use:   "status [FILE_HASH]",
		Short: "Get status of torrent(s) in progress.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileHash := ""
			if len(args) > 0 {
				fileHash = args[0]
			}

			conn, err := grpc.Dial(*daemonAddr, grpc.WithInsecure())
			if err != nil {
				return err
			}

			request := &daemon.StatusRequest{
				FileHash: fileHash,
			}

			client := daemon.NewDaemonClient(conn)
			response, err := client.Status(context.Background(), request)
			if err != nil {
				return err
			}

			fmt.Println(response.Torrents)

			return conn.Close()
		},
		Args: cobra.MaximumNArgs(1),
	}

	trackerAddr *string
	pieceSize   *int64
	daemonAddr  *string
	withData    *bool
)

func main() {
	trackerAddr = cmdCreate.Flags().StringP("tracker-address", "t", tracker.ListenAddress, "Tracker address.")
	pieceSize = cmdCreate.Flags().Int64P("piece-size", "s", piece.Size, "Piece size.")
	daemonAddr = cmdAdd.Flags().StringP("listen-address", "l", daemon.ListenAddress, "Daemon listening address.")
	withData = cmdDelete.Flags().BoolP("with-data", "d", false, "Delete also downloaded data from disk.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})
	cmdRoot.AddCommand(
		cmdAdd,
		cmdCreate,
		cmdDelete,
		cmdStatus,
	)

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
