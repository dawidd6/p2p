package seed

import "context"

type Seed struct {
	UnimplementedSeedServer
}

func NewSeed() *Seed {
	return &Seed{}
}

func (seed *Seed) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	return &GetResponse{}, nil
}
