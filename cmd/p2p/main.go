package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dawidd6/p2p/pkg/config"

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

			t, err := torrent.Create(file, *pieceSize, *trackerAddress)
			if err != nil {
				return err
			}

			err = torrent.Verify(t, file)
			if err != nil {
				return err
			}

			err = torrent.Save(t, "", t.FileName)
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

			address := net.JoinHostPort(*daemonHost, *daemonPort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
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

			address := net.JoinHostPort(*daemonHost, *daemonPort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
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

	cmdPause = &cobra.Command{
		Use:   "pause FILE_HASH",
		Short: "Pause specified torrent by file hash.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileHash := args[0]

			address := net.JoinHostPort(*daemonHost, *daemonPort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}

			request := &daemon.PauseRequest{
				FileHash: fileHash,
			}

			client := daemon.NewDaemonClient(conn)
			_, err = client.Pause(context.Background(), request)
			if err != nil {
				return err
			}

			return conn.Close()
		},
		Args: cobra.ExactArgs(1),
	}

	cmdResume = &cobra.Command{
		Use:   "resume FILE_HASH",
		Short: "Resume specified torrent by file hash.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileHash := args[0]

			address := net.JoinHostPort(*daemonHost, *daemonPort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}

			request := &daemon.ResumeRequest{
				FileHash: fileHash,
			}

			client := daemon.NewDaemonClient(conn)
			_, err = client.Resume(context.Background(), request)
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

			address := net.JoinHostPort(*daemonHost, *daemonPort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
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

			fmt.Printf("%-64v  %-6v  %-9v  %-16v  %-16v  %-5v\n",
				"FILE_HASH",
				"PAUSED",
				"COMPLETED",
				"DOWNLOADED_BYTES",
				"UPLOADED_BYTES",
				"RATIO",
			)
			for _, torrentState := range response.States {
				fmt.Printf("%-32v %-64v  %-6v  %-9v  %-16v  %-16v  %-16v  %-5.2f\n",
					torrentState.FileName,
					torrentState.FileHash,
					torrentState.Paused,
					torrentState.Completed,
					torrentState.PeersCount,
					torrentState.DownloadedBytes,
					torrentState.UploadedBytes,
					float32(torrentState.UploadedBytes)/float32(torrentState.DownloadedBytes),
				)
			}

			return conn.Close()
		},
		Args: cobra.MaximumNArgs(1),
	}

	daemonHost     *string
	daemonPort     *string
	trackerAddress *string
	pieceSize      *int64
	withData       *bool
)

func main() {
	conf := config.Default()

	daemonHost = cmdRoot.PersistentFlags().StringP("host", "l", conf.DaemonHost, "Daemon listening host.")
	daemonPort = cmdRoot.PersistentFlags().StringP("port", "p", conf.DaemonPort, "Daemon listening port.")
	trackerAddress = cmdCreate.Flags().StringP("tracker-address", "t", net.JoinHostPort(conf.TrackerHost, conf.TrackerPort), "Tracker address.")
	pieceSize = cmdCreate.Flags().Int64P("piece-size", "s", conf.PieceSize, "Piece size.")
	withData = cmdDelete.Flags().BoolP("with-data", "d", false, "Delete also downloaded data from disk.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})
	cmdRoot.AddCommand(
		cmdAdd,
		cmdCreate,
		cmdDelete,
		cmdPause,
		cmdResume,
		cmdStatus,
	)

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
