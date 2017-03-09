//go:generate mkdir -p pb
//go:generate protoc -I "../proto/" "../proto/service.proto" --go_out=plugins=grpc:pb
package client

import (
	"crypto/tls"
	"os"

	"./pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Connect(host string) (pb.TestServiceClient, *grpc.ClientConn, error) {
	// creds, err := credentials.NewClientTLSFromFile("../server.crt", "")
	// if err != nil {
	// 	return nil, nil, err
	// }

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	sslKeyLogfile := os.Getenv("SSLKEYLOGFILE")
	if sslKeyLogfile != "" {
		var w *os.File
		var err error
		w, err = os.OpenFile(sslKeyLogfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return nil, nil, err
		}
		tlsConfig.KeyLogWriter = w
	}
	creds := credentials.NewTLS(tlsConfig)
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewTestServiceClient(conn)

	return client, conn, nil
}
