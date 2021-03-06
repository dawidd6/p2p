// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package daemon

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// DaemonClient is the client API for Daemon service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DaemonClient interface {
	Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	Resume(ctx context.Context, in *ResumeRequest, opts ...grpc.CallOption) (*ResumeResponse, error)
	Pause(ctx context.Context, in *PauseRequest, opts ...grpc.CallOption) (*PauseResponse, error)
}

type daemonClient struct {
	cc grpc.ClientConnInterface
}

func NewDaemonClient(cc grpc.ClientConnInterface) DaemonClient {
	return &daemonClient{cc}
}

func (c *daemonClient) Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error) {
	out := new(AddResponse)
	err := c.cc.Invoke(ctx, "/Daemon/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/Daemon/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/Daemon/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) Resume(ctx context.Context, in *ResumeRequest, opts ...grpc.CallOption) (*ResumeResponse, error) {
	out := new(ResumeResponse)
	err := c.cc.Invoke(ctx, "/Daemon/Resume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) Pause(ctx context.Context, in *PauseRequest, opts ...grpc.CallOption) (*PauseResponse, error) {
	out := new(PauseResponse)
	err := c.cc.Invoke(ctx, "/Daemon/Pause", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DaemonServer is the server API for Daemon service.
// All implementations must embed UnimplementedDaemonServer
// for forward compatibility
type DaemonServer interface {
	Add(context.Context, *AddRequest) (*AddResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	Resume(context.Context, *ResumeRequest) (*ResumeResponse, error)
	Pause(context.Context, *PauseRequest) (*PauseResponse, error)
	mustEmbedUnimplementedDaemonServer()
}

// UnimplementedDaemonServer must be embedded to have forward compatible implementations.
type UnimplementedDaemonServer struct {
}

func (UnimplementedDaemonServer) Add(context.Context, *AddRequest) (*AddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedDaemonServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDaemonServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedDaemonServer) Resume(context.Context, *ResumeRequest) (*ResumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Resume not implemented")
}
func (UnimplementedDaemonServer) Pause(context.Context, *PauseRequest) (*PauseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pause not implemented")
}
func (UnimplementedDaemonServer) mustEmbedUnimplementedDaemonServer() {}

// UnsafeDaemonServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DaemonServer will
// result in compilation errors.
type UnsafeDaemonServer interface {
	mustEmbedUnimplementedDaemonServer()
}

func RegisterDaemonServer(s grpc.ServiceRegistrar, srv DaemonServer) {
	s.RegisterService(&_Daemon_serviceDesc, srv)
}

func _Daemon_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Daemon/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).Add(ctx, req.(*AddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Daemon/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Daemon/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_Resume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).Resume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Daemon/Resume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).Resume(ctx, req.(*ResumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_Pause_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PauseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).Pause(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Daemon/Pause",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).Pause(ctx, req.(*PauseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Daemon_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Daemon",
	HandlerType: (*DaemonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _Daemon_Add_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Daemon_Delete_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _Daemon_Status_Handler,
		},
		{
			MethodName: "Resume",
			Handler:    _Daemon_Resume_Handler,
		},
		{
			MethodName: "Pause",
			Handler:    _Daemon_Pause_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/daemon/daemon.proto",
}

// SeederClient is the client API for Seeder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SeederClient interface {
	Seed(ctx context.Context, in *SeedRequest, opts ...grpc.CallOption) (*SeedResponse, error)
}

type seederClient struct {
	cc grpc.ClientConnInterface
}

func NewSeederClient(cc grpc.ClientConnInterface) SeederClient {
	return &seederClient{cc}
}

func (c *seederClient) Seed(ctx context.Context, in *SeedRequest, opts ...grpc.CallOption) (*SeedResponse, error) {
	out := new(SeedResponse)
	err := c.cc.Invoke(ctx, "/Seeder/Seed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SeederServer is the server API for Seeder service.
// All implementations must embed UnimplementedSeederServer
// for forward compatibility
type SeederServer interface {
	Seed(context.Context, *SeedRequest) (*SeedResponse, error)
	mustEmbedUnimplementedSeederServer()
}

// UnimplementedSeederServer must be embedded to have forward compatible implementations.
type UnimplementedSeederServer struct {
}

func (UnimplementedSeederServer) Seed(context.Context, *SeedRequest) (*SeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Seed not implemented")
}
func (UnimplementedSeederServer) mustEmbedUnimplementedSeederServer() {}

// UnsafeSeederServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SeederServer will
// result in compilation errors.
type UnsafeSeederServer interface {
	mustEmbedUnimplementedSeederServer()
}

func RegisterSeederServer(s grpc.ServiceRegistrar, srv SeederServer) {
	s.RegisterService(&_Seeder_serviceDesc, srv)
}

func _Seeder_Seed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeederServer).Seed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Seeder/Seed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeederServer).Seed(ctx, req.(*SeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Seeder_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Seeder",
	HandlerType: (*SeederServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Seed",
			Handler:    _Seeder_Seed_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/daemon/daemon.proto",
}
