package seed

import "context"

const DefaultListenAddr = "0.0.0.0:44444"

type Seed struct {
	UnimplementedSeedServer
}

func NewSeed() *Seed {
	return &Seed{}
}

func (seed *Seed) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	return &GetResponse{}, nil
}
