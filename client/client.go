//go:generate mkdir -p pb
//go:generate protoc -I "../proto/" "../proto/service.proto" --go_out=plugins=grpc:pb
package client

import "./pb"
import "google.golang.org/grpc"
import "google.golang.org/grpc/credentials"

func Connect(host string) (pb.TestServiceClient, *grpc.ClientConn, error) {
	creds, err := credentials.NewClientTLSFromFile("../server.crt", "")
	if err != nil {
		return nil, nil, err
	}
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewTestServiceClient(conn)

	return client, conn, nil
}
