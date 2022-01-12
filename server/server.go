package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"w-r-grpc/pb"
	"w-r-grpc/platform"
	"w-r-grpc/platform/db"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedSessionServiceServer
}

func (*server) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	idStr := req.Id

	book := db.Select(idStr)

	resp := &pb.GetBookResponse{
		Id:        book.Id,
		Name:      book.Name,
		Author:    book.Author,
		Isbn:      book.Isbn,
		Publisher: book.Publisher,
		Genre:     book.Genre,
		Year:      book.YearOfPublication,
		Pages:     book.Pages,
	}

	return resp, nil
}

func (*server) PostBook(ctx context.Context, req *pb.PostBookRequest) (*pb.PostBookResponse, error) {
	strId := uuid.New().String()

	book := platform.Book{
		Id:                strId,
		Name:              req.Name,
		Author:            req.Author,
		Isbn:              req.Isbn,
		Publisher:         req.Publisher,
		Genre:             req.Genre,
		YearOfPublication: req.Year,
		Pages:             req.Pages,
	}

	db.Insert(book)

	resp := &pb.PostBookResponse{
		Result: "OK",
	}

	return resp, nil
}

func (*server) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	idStr := req.Id

	if !db.Delete(idStr) {
		return nil, fmt.Errorf("error while delete book")
	}

	return &pb.DeleteBookResponse{
		Result: "OK",
	}, nil

}

func (*server) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	idStr := req.Id

	book := db.Select(idStr)

	if book.Id != "" {

		book.Name = req.Name
		book.Author = req.Author
		book.Isbn = req.Isbn
		book.Publisher = req.Publisher
		book.Genre = req.Genre
		book.YearOfPublication = req.Year
		book.Pages = req.Pages

		db.Update(book)

		return &pb.UpdateBookResponse{
			Result: "OK",
		}, nil
	}

	return nil, fmt.Errorf("invalid id from request")
}

func StartGRPC(address string) {

	db.OpenDB()
	defer db.CloseDB()

	log.Println("grpc server is starting...")

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("error while listen tcp: %v", err)
	}

	defer lis.Close()

	s := grpc.NewServer() // empty options!!!! (for security should use tls)

	pb.RegisterSessionServiceServer(s, &server{})

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
