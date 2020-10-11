package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dawidd6/p2p/pkg/proto"
	"github.com/dawidd6/p2p/pkg/version"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	cmdRoot = &cobra.Command{
		Use:           "p2p-tracker",
		Short:         "P2P file sharing system based on gRPC. (tracker client)",
		Version:       version.Version,
		RunE:          run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	trackerAddr *string
)

func run(cmd *cobra.Command, args []string) error {
	conn, err := grpc.Dial(*trackerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	request := &proto.ListRequest{}
	client := proto.NewTrackerClient(conn)
	reply, err := client.List(context.TODO(), request)
	if err != nil {
		return err
	}

	fmt.Println(reply)

	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	trackerAddr = cmdRoot.Flags().StringP("address", "a", "localhost:8889", "Tracker daemon address.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
