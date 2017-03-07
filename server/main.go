//go:generate mkdir -p server/pb
//go:generate protoc -I "../proto/" "../proto/service.proto" --go_out=plugins=grpc:server/pb
package main

import "./server"
import "os"
import "fmt"

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: server <listen_host>\n")
		fmt.Printf("Example: server 127.0.0.1:4242\n")
		os.Exit(-1)
	}
	server.StartServer(os.Args[1])
}
