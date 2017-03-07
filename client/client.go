//go:generate mkdir -p pb
//go:generate protoc -I "../proto/" "../proto/service.proto" --go_out=plugins=grpc:pb
package client

import "./pb"
import "google.golang.org/grpc"

func Connect(host string) (pb.TestServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewTestServiceClient(conn)

	return client, conn, nil
}
