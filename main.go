package main

import "w-r-grpc/server"

const address string = "0.0.0.0:8080"

func main() {
	server.StartGRPC(address)
}
