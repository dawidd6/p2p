package main

import (
	"context"
	"log"
	"net"

	"github.com/dawidd6/p2p/pkg/proto"
	"github.com/dawidd6/p2p/pkg/version"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var cmdRoot = &cobra.Command{
	Use:           "p2p-trackerd",
	Short:         "P2P tracker daemon.",
	Version:       version.Version,
	RunE:          run,
	SilenceErrors: false,
	SilenceUsage:  false,
}

func run(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func sayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Printf("Received: %v", in.Message)
	return &proto.HelloResponse{Message: "Hello " + in.Message}, nil
}

func main() {
	err := cmdRoot.Execute()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	proto.RegisterTrackerServiceService(server, &proto.TrackerServiceService{SayHello: sayHello})
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
