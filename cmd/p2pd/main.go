package main

import (
	"log"
	"net"

	"github.com/dawidd6/p2p/pkg/daemon"
	"github.com/dawidd6/p2p/pkg/proto"
	"github.com/dawidd6/p2p/pkg/version"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	cmdRoot = &cobra.Command{
		Use:           "p2pd",
		Short:         "P2P file sharing system based on gRPC (daemon).",
		Version:       version.Version,
		RunE:          run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	listenAddr *string
)

func run(cmd *cobra.Command, args []string) error {
	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		return err
	}

	d := daemon.NewDaemon(*listenAddr)
	g := grpc.NewServer()

	proto.RegisterDaemonServer(g, d)

	return g.Serve(listener)
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	listenAddr = cmdRoot.Flags().StringP("address", "a", "localhost:8888", "Address on which daemon should listen.")

	cmdRoot.SetHelpCommand(&cobra.Command{Hidden: true})

	err := cmdRoot.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
