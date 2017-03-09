# caddyGrpcTest
A little test for GRPC proxying with [Caddy](https://caddyserver.com/)

# Goal
This is supposed to easily test the functionality of reverse-proxying GRPC-Services with caddy.

# What's in this?
This contains a grpc server and a grpc client, both written in go.
The client implements one test for each kind of grpc call (see [grpc-concepts](http://www.grpc.io/docs/guides/concepts.html) for more details).
Using the go testsuite the client will tell you what is working and what is not.

I have also provided a `Caddyfile` as a sort of common ground for testing. This can of course be modified.

# Requirements
You need the following things:

- Go v1.8
- protobuf (v3)
- grpc
- Caddy

## Setup of protobuf and grpc
At this point I assume you have a complete and working go-setup.

1. Download the protobuf compiler from [github](https://github.com/google/protobuf/releases).
Any version 3.x should work, but I would advise against the alpha builds.
*Note:* You only need protoc-3.x-os-arch.zip
2. Copy the "protoc" executable to a folder included in your `$PATH`.
3. Install the go protobuf libraries. `$ go get -u github.com/golang/protobuf/protoc-gen-go`
4. Make sure that your `$GOPATH/bin` is part of your `$PATH` as that's where the "protoc-gen-go" executale will be installed.
5. Install the go grpc libraries. `$ go get google.golang.org/grpc`

## Generate the necessary code
Assuming you are currently in the root of this repo.

```
$ cd client
$ go generate
$ cd ../server
$ go generate
```

# Running the tests
The tests can run on multiple servers.
In order for this to work you are required to set the environment variable `TEST_HOSTS` in the format `host1:port1,host2:port2,...`.
This allowes you to test multiple Caddy configurations at once, as well as being able to test the unproxied server.
E. g. `TEST_HOSTS=127.0.0.1:443,127.0.0.1:4242`.

So to run the tests do

1. Start the grpc server. `$ go run server/main.go ":<port>" &` where `<port>` is the port you want to run it on.
2. Start the caddy server. `$ caddy -conf Caddyfile &`
3. Specify the servers under test. `$ export TEST_HOSTS="..."`
4. Run the tests. `$ cd client && go test -v`

Grpc will actually write some errors to the `stderr` on its own. So if you just want to see what works and what doesn't just run `go test -v 2>/dev/null` instead. I would advise the `-v` flag though because otherwise you won't see successfull tests listed.

If you want to sniff the network traffic you can `export SSLKEYLOGFILE="/path/to/file"` and use that file in wireshark in order to decode the encrypted data.

# License

This is released into the wild under the MIT license, so basically do what you will with it. :wink:
