package client

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"./pb"
)

var connHosts []string

func init() {
	connHosts = make([]string, 0)

	hosts := os.Getenv("TEST_HOSTS")
	if hosts == "" {
		fmt.Println("You did not specify any hosts.\nDo so by setting the environment variable \"TEST_HOSTS\"=\"host1:port1,host2:port2,...\".")
		os.Exit(-1)
	} else {
		connHosts = strings.Split(hosts, ",")
	}
}

func TestSimpleEcho(t *testing.T) {
	for _, host := range connHosts {
		t.Run(host, func(t *testing.T) {
			client, conn, err := Connect(host)
			if err != nil {
				t.Error(err)
				return
			}
			defer conn.Close()

			reqMsg := "This is an echo Test."
			expected := reqMsg
			resp, err := client.SimpleEcho(context.Background(), &pb.EchoRequest{reqMsg, 1})
			if err != nil {
				t.Log("Error while sending request.")
				t.Error(err)
			} else if resp.Message != expected {
				t.Errorf("Invalid response. Expected \"%s\". Got \"%s\"", expected, resp.Message)
			}
		})
	}
}

func TestServerStreamEcho(t *testing.T) {
	for _, host := range connHosts {
		t.Run(host, func(t *testing.T) {
			client, conn, err := Connect(host)
			if err != nil {
				t.Error(err)
				return
			}
			defer conn.Close()

			reqMsg := "This is an echo Test."
			expectedMsg := reqMsg
			reqCount := int32(6)
			expectedCount := reqCount
			stream, err := client.ServerStreamEcho(context.Background(), &pb.EchoRequest{reqMsg, reqCount})
			if err != nil {
				t.Log("Error while creating stream.")
				t.Error(err)
				return
			}
			numRecvdPackages := int32(0)
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					break
				}
				numRecvdPackages++
				if err != nil {
					t.Log("Error while reading stream.")
					t.Error(err)
					break
				} else if resp.Message != expectedMsg {
					t.Errorf("Invalid response on the %d-th package. Expected \"%s\". Got \"%s\"", numRecvdPackages, expectedMsg, resp.Message)
				}
			}
			if expectedCount != numRecvdPackages {
				t.Errorf("Expected %d packages. Got %d.", expectedCount, numRecvdPackages)
			}
		})
	}
}

func TestClientStreamEcho(t *testing.T) {
	for _, host := range connHosts {
		t.Run(host, func(t *testing.T) {
			client, conn, err := Connect(host)
			if err != nil {
				t.Error(err)
				return
			}
			defer conn.Close()

			reqMsg := "This is an echo Test."
			reqCount := int32(6)
			expectedMsg := ""
			for i := int32(0); i < reqCount; i++ {
				expectedMsg += reqMsg
			}
			stream, err := client.ClientStreamEcho(context.Background())
			if err != nil {
				t.Log("Error while creating stream.")
				t.Error(err)
				return
			}
			for i := int32(0); i < reqCount; i++ {
				if err := stream.Send(&pb.EchoRequest{reqMsg, 1}); err != nil {
					t.Log("Error while writing to stream.")
					t.Error(err)
					break
				}
			}
			resp, err := stream.CloseAndRecv()
			if err != nil {
				t.Log("Error while receiving response.")
				t.Error(err)
			} else if resp.Message != expectedMsg {
				t.Errorf("Invalid response. Expected \"%s\". Got \"%s\"", expectedMsg, resp.Message)
			}
		})
	}
}

func TestBidirectionalStreamEcho(t *testing.T) {
	for _, host := range connHosts {
		t.Run(host, func(t *testing.T) {
			client, conn, err := Connect(host)
			if err != nil {
				t.Error(err)
				return
			}
			defer conn.Close()

			reqMsg := "This is an echo Test."
			expectedMsg := reqMsg
			reqCount := int32(8)
			expectedCount := reqCount
			stream, err := client.BidirectionalStreamEcho(context.Background())
			if err != nil {
				t.Log("Error while creating stream.")
				t.Error(err)
				return
			}
			numRecvdPackages := int32(0)
			waitc := make(chan struct{})
			go func() {
				for {
					resp, err := stream.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						if !t.Failed() {
							t.Log("Error while reading from stream.")
							t.Error(err)
						}
						break
					}
					numRecvdPackages++
					if resp.Message != expectedMsg {
						t.Errorf("Invalid response on the %d-th package. Expected \"%s\". Got \"%s\"", numRecvdPackages, expectedMsg, resp.Message)
						break
					}
				}
				close(waitc)
			}()
			for i := int32(0); i < reqCount; i++ {
				if err := stream.Send(&pb.EchoRequest{reqMsg, 1}); err != nil {
					t.Log("Error while writing to stream.")
					t.Error(err)
					break
				}
			}
			if err := stream.CloseSend(); !t.Failed() && err != nil {
				t.Log("Error while closing stream.")
				t.Error(err)
			}
			<-waitc
			if !t.Failed() && expectedCount != numRecvdPackages {
				t.Errorf("Expected %d packages. Got %d.", expectedCount, numRecvdPackages)
			}
		})
	}
}

// func Test____(t *testing.T) {
// 	for _, host := range connHosts {
// 		t.Run(host, func(t *testing.T) {
// 			client, conn, err := Connect(host)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			defer conn.Close()
//
// 			// CODE HERE
// 		})
// 	}
// }
