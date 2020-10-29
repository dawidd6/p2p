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
			dataFilePath := args[0]

			dataFile, err := os.OpenFile(dataFilePath, os.O_RDONLY, 0666)
			if err != nil {
				return err
			}

			trackerAddress := net.JoinHostPort(conf.TrackerHost, conf.TrackerPort)
			torr, err := torrent.Create(dataFile, conf.PieceSize, trackerAddress)
			if err != nil {
				return err
			}

			torrentFilePath := torrent.File("", torr.FileName)
			torrentFile, err := os.OpenFile(torrentFilePath, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				return err
			}

			err = torrent.Verify(dataFile, torr)
			if err != nil {
				return err
			}

			err = torrent.Write(torrentFile, torr)
			if err != nil {
				return err
			}

			err = dataFile.Close()
			if err != nil {
				return err
			}

			err = torrentFile.Close()
			if err != nil {
				return err
			}

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	cmdAdd = &cobra.Command{
		Use:   "add TORRENT_FILE",
		Short: "Add specified torrents to client.",
		RunE: func(cmd *cobra.Command, args []string) error {
			torrentFilePath := args[0]

			torrentFile, err := os.OpenFile(torrentFilePath, os.O_RDONLY, 0666)
			if err != nil {
				return err
			}

			torr, err := torrent.Read(torrentFile)
			if err != nil {
				return err
			}

			err = torrentFile.Close()
			if err != nil {
				return err
			}

			address := net.JoinHostPort(conf.DaemonHost, conf.DaemonPort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}

			request := &daemon.AddRequest{
				Torrent: torr,
			}

			ctx, cancel := context.WithTimeout(context.Background(), conf.CallTimeout)
			defer cancel()

			client := daemon.NewDaemonClient(conn)
			_, err = client.Add(ctx, request)
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

			address := net.JoinHostPort(conf.DaemonHost, conf.DaemonPort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}

			request := &daemon.DeleteRequest{
				FileHash: fileHash,
				WithData: conf.DeleteWithData,
			}

			ctx, cancel := context.WithTimeout(context.Background(), conf.CallTimeout)
			defer cancel()

			client := daemon.NewDaemonClient(conn)
			_, err = client.Delete(ctx, request)
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

			address := net.JoinHostPort(conf.DaemonHost, conf.DaemonPort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}

			request := &daemon.PauseRequest{
				FileHash: fileHash,
			}

			ctx, cancel := context.WithTimeout(context.Background(), conf.CallTimeout)
			defer cancel()

			client := daemon.NewDaemonClient(conn)
			_, err = client.Pause(ctx, request)
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

			address := net.JoinHostPort(conf.DaemonHost, conf.DaemonPort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}

			request := &daemon.ResumeRequest{
				FileHash: fileHash,
			}

			ctx, cancel := context.WithTimeout(context.Background(), conf.CallTimeout)
			defer cancel()

			client := daemon.NewDaemonClient(conn)
			_, err = client.Resume(ctx, request)
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

			address := net.JoinHostPort(conf.DaemonHost, conf.DaemonPort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return err
			}

			request := &daemon.StatusRequest{
				FileHash: fileHash,
			}

			ctx, cancel := context.WithTimeout(context.Background(), conf.CallTimeout)
			defer cancel()

			client := daemon.NewDaemonClient(conn)
			response, err := client.Status(ctx, request)
			if err != nil {
				return err
			}

			fmt.Printf("%-32v  %-64v  %-6v  %-9v  %-16v  %-16v  %-5v\n",
				"FILE_NAME",
				"FILE_HASH",
				"PAUSED",
				"COMPLETED",
				"DOWNLOADED_BYTES",
				"UPLOADED_BYTES",
				"RATIO",
			)
			for _, torrentState := range response.States {
				fmt.Printf("%-32v  %-64v  %-6v  %-9v  %-16v  %-16v  %-16v  %-5.2f\n",
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

	conf = config.Default()
)

func main() {
	cmdRoot.PersistentFlags().StringVar(&conf.DaemonHost, "host", conf.DaemonHost, "Daemon listening host.")
	cmdRoot.PersistentFlags().StringVar(&conf.DaemonPort, "port", conf.DaemonPort, "Daemon listening port.")
	cmdRoot.PersistentFlags().DurationVar(&conf.CallTimeout, "call-timeout", conf.CallTimeout, "Delete also downloaded data from disk.")

	cmdCreate.Flags().StringVar(&conf.TrackerHost, "tracker-host", conf.TrackerHost, "Tracker address host.")
	cmdCreate.Flags().StringVar(&conf.TrackerPort, "tracker-port", conf.TrackerPort, "Tracker address port.")
	cmdCreate.Flags().Int64Var(&conf.PieceSize, "piece-size", conf.PieceSize, "Piece size.")

	cmdDelete.Flags().BoolVar(&conf.DeleteWithData, "with-data", conf.DeleteWithData, "Delete also downloaded data from disk.")

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
