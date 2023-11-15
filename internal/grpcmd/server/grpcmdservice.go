package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/grpcmd/grpcmd/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedGrpcmdServiceServer
}

func (server) UnaryMethod(_ context.Context, req *pb.GrpcmdRequest) (*pb.GrpcmdResponse, error) {
	return &pb.GrpcmdResponse{Message: "Hello, " + req.Name + "!"}, nil
}

func (server) ClientStreamingMethod(stream pb.GrpcmdService_ClientStreamingMethodServer) error {
	names := make([]string, 0, 5)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		names = append(names, req.Name)
	}
	stream.SendAndClose(&pb.GrpcmdResponse{Message: "Hello, " + strings.Join(names, " + ") + "!"})
	return nil
}

func (server) ServerStreamingMethod(req *pb.GrpcmdRequest, stream pb.GrpcmdService_ServerStreamingMethodServer) error {
	stream.Send(&pb.GrpcmdResponse{Message: "Hello, "})
	for _, v := range req.Name {
		stream.Send(&pb.GrpcmdResponse{Message: string(v)})
	}
	stream.Send(&pb.GrpcmdResponse{Message: "!"})
	return nil
}

func (server) BidirectionalStreamingMethod(stream pb.GrpcmdService_BidirectionalStreamingMethodServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		stream.Send(&pb.GrpcmdResponse{Message: "Hello, " + req.Name + "!"})
	}
}

func Run(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterGrpcmdServiceServer(s, &server{})
	reflection.Register(s)

	fmt.Printf("Listening on address: %s\n\n", address)
	fmt.Println("Try running:\n\tgrpc " + address + " UnaryMethod '{\"name\": \"Bob\"}'")

	err = s.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
