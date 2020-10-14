package seed

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

const DefaultListenAddr = "0.0.0.0:44444"

type Seed struct {
	address string
	UnimplementedSeedServer
}

func New(address string) *Seed {
	return &Seed{
		address: address,
	}
}

func (seed *Seed) Listen() error {
	listener, err := net.Listen("tcp", seed.address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterSeedServer(server, seed)
	return server.Serve(listener)
}

func (seed *Seed) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	return &GetResponse{}, nil
}
