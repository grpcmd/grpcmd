// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: proto/grpcmd_service.proto

package pb

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

const (
	GrpcmdService_UnaryMethod_FullMethodName                  = "/grpcmd.GrpcmdService/UnaryMethod"
	GrpcmdService_ClientStreamingMethod_FullMethodName        = "/grpcmd.GrpcmdService/ClientStreamingMethod"
	GrpcmdService_ServerStreamingMethod_FullMethodName        = "/grpcmd.GrpcmdService/ServerStreamingMethod"
	GrpcmdService_BidirectionalStreamingMethod_FullMethodName = "/grpcmd.GrpcmdService/BidirectionalStreamingMethod"
)

// GrpcmdServiceClient is the client API for GrpcmdService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GrpcmdServiceClient interface {
	UnaryMethod(ctx context.Context, in *GrpcmdRequest, opts ...grpc.CallOption) (*GrpcmdResponse, error)
	ClientStreamingMethod(ctx context.Context, opts ...grpc.CallOption) (GrpcmdService_ClientStreamingMethodClient, error)
	ServerStreamingMethod(ctx context.Context, in *GrpcmdRequest, opts ...grpc.CallOption) (GrpcmdService_ServerStreamingMethodClient, error)
	BidirectionalStreamingMethod(ctx context.Context, opts ...grpc.CallOption) (GrpcmdService_BidirectionalStreamingMethodClient, error)
}

type grpcmdServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGrpcmdServiceClient(cc grpc.ClientConnInterface) GrpcmdServiceClient {
	return &grpcmdServiceClient{cc}
}

func (c *grpcmdServiceClient) UnaryMethod(ctx context.Context, in *GrpcmdRequest, opts ...grpc.CallOption) (*GrpcmdResponse, error) {
	out := new(GrpcmdResponse)
	err := c.cc.Invoke(ctx, GrpcmdService_UnaryMethod_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcmdServiceClient) ClientStreamingMethod(ctx context.Context, opts ...grpc.CallOption) (GrpcmdService_ClientStreamingMethodClient, error) {
	stream, err := c.cc.NewStream(ctx, &GrpcmdService_ServiceDesc.Streams[0], GrpcmdService_ClientStreamingMethod_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcmdServiceClientStreamingMethodClient{stream}
	return x, nil
}

type GrpcmdService_ClientStreamingMethodClient interface {
	Send(*GrpcmdRequest) error
	CloseAndRecv() (*GrpcmdResponse, error)
	grpc.ClientStream
}

type grpcmdServiceClientStreamingMethodClient struct {
	grpc.ClientStream
}

func (x *grpcmdServiceClientStreamingMethodClient) Send(m *GrpcmdRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *grpcmdServiceClientStreamingMethodClient) CloseAndRecv() (*GrpcmdResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(GrpcmdResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *grpcmdServiceClient) ServerStreamingMethod(ctx context.Context, in *GrpcmdRequest, opts ...grpc.CallOption) (GrpcmdService_ServerStreamingMethodClient, error) {
	stream, err := c.cc.NewStream(ctx, &GrpcmdService_ServiceDesc.Streams[1], GrpcmdService_ServerStreamingMethod_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcmdServiceServerStreamingMethodClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GrpcmdService_ServerStreamingMethodClient interface {
	Recv() (*GrpcmdResponse, error)
	grpc.ClientStream
}

type grpcmdServiceServerStreamingMethodClient struct {
	grpc.ClientStream
}

func (x *grpcmdServiceServerStreamingMethodClient) Recv() (*GrpcmdResponse, error) {
	m := new(GrpcmdResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *grpcmdServiceClient) BidirectionalStreamingMethod(ctx context.Context, opts ...grpc.CallOption) (GrpcmdService_BidirectionalStreamingMethodClient, error) {
	stream, err := c.cc.NewStream(ctx, &GrpcmdService_ServiceDesc.Streams[2], GrpcmdService_BidirectionalStreamingMethod_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcmdServiceBidirectionalStreamingMethodClient{stream}
	return x, nil
}

type GrpcmdService_BidirectionalStreamingMethodClient interface {
	Send(*GrpcmdRequest) error
	Recv() (*GrpcmdResponse, error)
	grpc.ClientStream
}

type grpcmdServiceBidirectionalStreamingMethodClient struct {
	grpc.ClientStream
}

func (x *grpcmdServiceBidirectionalStreamingMethodClient) Send(m *GrpcmdRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *grpcmdServiceBidirectionalStreamingMethodClient) Recv() (*GrpcmdResponse, error) {
	m := new(GrpcmdResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GrpcmdServiceServer is the server API for GrpcmdService service.
// All implementations must embed UnimplementedGrpcmdServiceServer
// for forward compatibility
type GrpcmdServiceServer interface {
	UnaryMethod(context.Context, *GrpcmdRequest) (*GrpcmdResponse, error)
	ClientStreamingMethod(GrpcmdService_ClientStreamingMethodServer) error
	ServerStreamingMethod(*GrpcmdRequest, GrpcmdService_ServerStreamingMethodServer) error
	BidirectionalStreamingMethod(GrpcmdService_BidirectionalStreamingMethodServer) error
	mustEmbedUnimplementedGrpcmdServiceServer()
}

// UnimplementedGrpcmdServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGrpcmdServiceServer struct {
}

func (UnimplementedGrpcmdServiceServer) UnaryMethod(context.Context, *GrpcmdRequest) (*GrpcmdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnaryMethod not implemented")
}
func (UnimplementedGrpcmdServiceServer) ClientStreamingMethod(GrpcmdService_ClientStreamingMethodServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStreamingMethod not implemented")
}
func (UnimplementedGrpcmdServiceServer) ServerStreamingMethod(*GrpcmdRequest, GrpcmdService_ServerStreamingMethodServer) error {
	return status.Errorf(codes.Unimplemented, "method ServerStreamingMethod not implemented")
}
func (UnimplementedGrpcmdServiceServer) BidirectionalStreamingMethod(GrpcmdService_BidirectionalStreamingMethodServer) error {
	return status.Errorf(codes.Unimplemented, "method BidirectionalStreamingMethod not implemented")
}
func (UnimplementedGrpcmdServiceServer) mustEmbedUnimplementedGrpcmdServiceServer() {}

// UnsafeGrpcmdServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GrpcmdServiceServer will
// result in compilation errors.
type UnsafeGrpcmdServiceServer interface {
	mustEmbedUnimplementedGrpcmdServiceServer()
}

func RegisterGrpcmdServiceServer(s grpc.ServiceRegistrar, srv GrpcmdServiceServer) {
	s.RegisterService(&GrpcmdService_ServiceDesc, srv)
}

func _GrpcmdService_UnaryMethod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcmdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcmdServiceServer).UnaryMethod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcmdService_UnaryMethod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcmdServiceServer).UnaryMethod(ctx, req.(*GrpcmdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcmdService_ClientStreamingMethod_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GrpcmdServiceServer).ClientStreamingMethod(&grpcmdServiceClientStreamingMethodServer{stream})
}

type GrpcmdService_ClientStreamingMethodServer interface {
	SendAndClose(*GrpcmdResponse) error
	Recv() (*GrpcmdRequest, error)
	grpc.ServerStream
}

type grpcmdServiceClientStreamingMethodServer struct {
	grpc.ServerStream
}

func (x *grpcmdServiceClientStreamingMethodServer) SendAndClose(m *GrpcmdResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *grpcmdServiceClientStreamingMethodServer) Recv() (*GrpcmdRequest, error) {
	m := new(GrpcmdRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GrpcmdService_ServerStreamingMethod_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GrpcmdRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GrpcmdServiceServer).ServerStreamingMethod(m, &grpcmdServiceServerStreamingMethodServer{stream})
}

type GrpcmdService_ServerStreamingMethodServer interface {
	Send(*GrpcmdResponse) error
	grpc.ServerStream
}

type grpcmdServiceServerStreamingMethodServer struct {
	grpc.ServerStream
}

func (x *grpcmdServiceServerStreamingMethodServer) Send(m *GrpcmdResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _GrpcmdService_BidirectionalStreamingMethod_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GrpcmdServiceServer).BidirectionalStreamingMethod(&grpcmdServiceBidirectionalStreamingMethodServer{stream})
}

type GrpcmdService_BidirectionalStreamingMethodServer interface {
	Send(*GrpcmdResponse) error
	Recv() (*GrpcmdRequest, error)
	grpc.ServerStream
}

type grpcmdServiceBidirectionalStreamingMethodServer struct {
	grpc.ServerStream
}

func (x *grpcmdServiceBidirectionalStreamingMethodServer) Send(m *GrpcmdResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *grpcmdServiceBidirectionalStreamingMethodServer) Recv() (*GrpcmdRequest, error) {
	m := new(GrpcmdRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GrpcmdService_ServiceDesc is the grpc.ServiceDesc for GrpcmdService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GrpcmdService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpcmd.GrpcmdService",
	HandlerType: (*GrpcmdServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UnaryMethod",
			Handler:    _GrpcmdService_UnaryMethod_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ClientStreamingMethod",
			Handler:       _GrpcmdService_ClientStreamingMethod_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "ServerStreamingMethod",
			Handler:       _GrpcmdService_ServerStreamingMethod_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BidirectionalStreamingMethod",
			Handler:       _GrpcmdService_BidirectionalStreamingMethod_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/grpcmd_service.proto",
}