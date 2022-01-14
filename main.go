package main

import (
	"os"
	"w-r-grpc/server"
)

func main() {
	var address string

	if len(os.Args) < 2 {
		address = "0.0.0.0:8080"
	} else {
		address = os.Args[1]
	}

	server.StartGRPC(address)
}
