// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: muapp/v1/hello.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MuAppServiceClient is the client API for MuAppService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MuAppServiceClient interface {
	// Hello is a simple RPC that returns a greeting.
	Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
}

type muAppServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMuAppServiceClient(cc grpc.ClientConnInterface) MuAppServiceClient {
	return &muAppServiceClient{cc}
}

func (c *muAppServiceClient) Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, "/muapp.v1.MuAppService/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MuAppServiceServer is the server API for MuAppService service.
// All implementations must embed UnimplementedMuAppServiceServer
// for forward compatibility
type MuAppServiceServer interface {
	// Hello is a simple RPC that returns a greeting.
	Hello(context.Context, *HelloRequest) (*HelloResponse, error)
	mustEmbedUnimplementedMuAppServiceServer()
}

// UnimplementedMuAppServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMuAppServiceServer struct {
}

func (UnimplementedMuAppServiceServer) Hello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedMuAppServiceServer) mustEmbedUnimplementedMuAppServiceServer() {}

// UnsafeMuAppServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MuAppServiceServer will
// result in compilation errors.
type UnsafeMuAppServiceServer interface {
	mustEmbedUnimplementedMuAppServiceServer()
}

func RegisterMuAppServiceServer(s grpc.ServiceRegistrar, srv MuAppServiceServer) {
	s.RegisterService(&MuAppService_ServiceDesc, srv)
}

func _MuAppService_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MuAppServiceServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/muapp.v1.MuAppService/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MuAppServiceServer).Hello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MuAppService_ServiceDesc is the grpc.ServiceDesc for MuAppService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MuAppService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "muapp.v1.MuAppService",
	HandlerType: (*MuAppServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _MuAppService_Hello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "muapp/v1/hello.proto",
}
