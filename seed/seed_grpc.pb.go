// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package seed

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SeedClient is the client API for Seed service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SeedClient interface {
	Seed(ctx context.Context, opts ...grpc.CallOption) (Seed_SeedClient, error)
}

type seedClient struct {
	cc grpc.ClientConnInterface
}

func NewSeedClient(cc grpc.ClientConnInterface) SeedClient {
	return &seedClient{cc}
}

var seedSeedStreamDesc = &grpc.StreamDesc{
	StreamName:    "Seed",
	ServerStreams: true,
	ClientStreams: true,
}

func (c *seedClient) Seed(ctx context.Context, opts ...grpc.CallOption) (Seed_SeedClient, error) {
	stream, err := c.cc.NewStream(ctx, seedSeedStreamDesc, "/Seed/Seed", opts...)
	if err != nil {
		return nil, err
	}
	x := &seedSeedClient{stream}
	return x, nil
}

type Seed_SeedClient interface {
	Send(*SeedRequest) error
	Recv() (*SeedResponse, error)
	grpc.ClientStream
}

type seedSeedClient struct {
	grpc.ClientStream
}

func (x *seedSeedClient) Send(m *SeedRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *seedSeedClient) Recv() (*SeedResponse, error) {
	m := new(SeedResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SeedService is the service API for Seed service.
// Fields should be assigned to their respective handler implementations only before
// RegisterSeedService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type SeedService struct {
	Seed func(Seed_SeedServer) error
}

func (s *SeedService) seed(_ interface{}, stream grpc.ServerStream) error {
	return s.Seed(&seedSeedServer{stream})
}

type Seed_SeedServer interface {
	Send(*SeedResponse) error
	Recv() (*SeedRequest, error)
	grpc.ServerStream
}

type seedSeedServer struct {
	grpc.ServerStream
}

func (x *seedSeedServer) Send(m *SeedResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *seedSeedServer) Recv() (*SeedRequest, error) {
	m := new(SeedRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RegisterSeedService registers a service implementation with a gRPC server.
func RegisterSeedService(s grpc.ServiceRegistrar, srv *SeedService) {
	srvCopy := *srv
	if srvCopy.Seed == nil {
		srvCopy.Seed = func(Seed_SeedServer) error {
			return status.Errorf(codes.Unimplemented, "method Seed not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "Seed",
		Methods:     []grpc.MethodDesc{},
		Streams: []grpc.StreamDesc{
			{
				StreamName:    "Seed",
				Handler:       srvCopy.seed,
				ServerStreams: true,
				ClientStreams: true,
			},
		},
		Metadata: "seed/seed.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewSeedService creates a new SeedService containing the
// implemented methods of the Seed service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewSeedService(s interface{}) *SeedService {
	ns := &SeedService{}
	if h, ok := s.(interface{ Seed(Seed_SeedServer) error }); ok {
		ns.Seed = h.Seed
	}
	return ns
}

// UnstableSeedService is the service API for Seed service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableSeedService interface {
	Seed(Seed_SeedServer) error
}
