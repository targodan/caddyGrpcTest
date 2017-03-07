package server

import (
	"io"
	"log"
	"net"

	"./pb"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

type testServiceServer struct{}

func (s *testServiceServer) SimpleEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{req.Message}, nil
}

func (s *testServiceServer) ServerStreamEcho(req *pb.EchoRequest, stream pb.TestService_ServerStreamEchoServer) error {
	for i := 0; i < int(req.Count); i++ {
		if err := stream.Send(&pb.EchoResponse{req.Message}); err != nil {
			return err
		}
	}
	return nil
}

func (s *testServiceServer) ClientStreamEcho(stream pb.TestService_ClientStreamEchoServer) error {
	msg := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.EchoResponse{msg})
		}
		if err != nil {
			return err
		}
		msg += req.Message
	}
}

func (s *testServiceServer) BidirectionalStreamEcho(stream pb.TestService_BidirectionalStreamEchoServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err := stream.Send(&pb.EchoResponse{req.Message}); err != nil {
			return err
		}
	}
}

func StartServer(host string) {
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTestServiceServer(grpcServer, &testServiceServer{})
	grpcServer.Serve(lis)
}
