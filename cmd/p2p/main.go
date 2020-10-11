package main

import (
	"context"

	"github.com/dawidd6/p2p/pkg/proto"
	"github.com/dawidd6/p2p/pkg/torrent"
	"github.com/dawidd6/p2p/pkg/version"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	cmdRoot = &cobra.Command{
		Use:           "p2p",
		Short:         "P2P file sharing system based on gRPC.",
		Version:       version.Version,
		RunE:          run,
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

	torrentName *string
	torrentDir  *string
	daemonAddr  *string
)

func init() {
	daemonAddr = cmdRoot.PersistentFlags().StringP("address", "a", "localhost:8888", "Daemon address.")
	torrentName = cmdCreate.Flags().StringP("name", "n", "", "Torrent name.")
	torrentDir = cmdCreate.Flags().StringP("dir", "d", "", "Where files should be stored.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})
	cmdRoot.AddCommand(
		cmdAdd,
		cmdCreate,
	)
}

func run(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

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

	request := &proto.AddRequest{Torrent: torr}
	client := proto.NewDaemonClient(conn)
	_, err = client.Add(context.TODO(), request)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := cmdRoot.Execute()
	if err != nil {
		panic(err)
	}
}
