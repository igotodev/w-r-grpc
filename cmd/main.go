package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"w-r-grpc/config"
	"w-r-grpc/internal/domain/dblogic"
	"w-r-grpc/internal/server"
	"w-r-grpc/pb"
	"w-r-grpc/pkg/postgres"
)

func main() {
	var address string

	if len(os.Args) < 2 {
		address = "0.0.0.0:8080"
	} else {
		address = os.Args[1]
	}

	log.Println("grpc server is starting...")

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("error while listen tcp: %v", err)
	}

	defer lis.Close()

	s := grpc.NewServer() // empty options!!!! (for security should use tls)

	cfg := config.InitConfig()
	psql, err := postgres.OpenDB(cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	dbl := dblogic.NewDatabase(psql)

	pb.RegisterSessionServiceServer(s, &server.Server{Storage: dbl})

	// register reflection
	reflection.Register(s)

	chWair := make(chan os.Signal, 1)

	signal.Notify(chWair, os.Interrupt)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("error while serve: %v", err)
			close(chWair)
		}
	}()

	<-chWair
	s.Stop()

	log.Println("grpc server is stoped")
}
