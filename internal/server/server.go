package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"w-r-grpc/internal/domain/bookvalidator"
	"w-r-grpc/internal/domain/entity"
	"w-r-grpc/internal/domain/storage"
	"w-r-grpc/pb"
)

type Server struct {
	pb.UnimplementedSessionServiceServer
	Storage storage.Storage
}

func (s *Server) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	idStr := req.Id

	book := s.Storage.GetOne(idStr)

	if book.Id == "" {
		return nil, fmt.Errorf("error: invalid id")
	}

	bookResp := pb.Book{
		Id:        book.Id,
		Name:      book.Name,
		Author:    book.Author,
		Isbn:      book.Isbn,
		Publisher: book.Publisher,
		Genre:     book.Genre,
		Year:      book.YearOfPublication,
		Pages:     book.Pages,
	}

	resp := &pb.GetBookResponse{
		Book: &bookResp,
	}

	return resp, nil
}

func (s *Server) GetAllBooks(req *pb.GetAllBooksRequest, stream pb.SessionService_GetAllBooksServer) error {
	books := s.Storage.GetAll()

	for _, book := range books {
		bookResp := pb.Book{
			Id:        book.Id,
			Name:      book.Name,
			Author:    book.Author,
			Isbn:      book.Isbn,
			Publisher: book.Publisher,
			Genre:     book.Genre,
			Year:      book.YearOfPublication,
			Pages:     book.Pages,
		}

		resp := &pb.GetAllBooksResponse{
			Book: &bookResp,
		}
		stream.Send(resp)
	}

	return nil
}

func (s *Server) PostBook(ctx context.Context, req *pb.PostBookRequest) (*pb.PostBookResponse, error) {
	strId := uuid.New().String()

	book := entity.Book{
		Id:                strId,
		Name:              req.Book.Name,
		Author:            req.Book.Author,
		Isbn:              req.Book.Isbn,
		Publisher:         req.Book.Publisher,
		Genre:             req.Book.Genre,
		YearOfPublication: req.Book.Year,
		Pages:             req.Book.Pages,
	}

	if !bookvalidator.IsValid(book) {
		return nil, fmt.Errorf("error: invalid data")
	}

	if err := s.Storage.Post(book); err != nil {
		return &pb.PostBookResponse{
			Result: "error",
		}, err
	}

	resp := &pb.PostBookResponse{
		Result: "OK",
	}

	return resp, nil
}

func (s *Server) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	idStr := req.Id

	if err := s.Storage.Delete(idStr); err != nil {
		return &pb.DeleteBookResponse{
			Result: "error",
		}, err
	}

	return &pb.DeleteBookResponse{
		Result: "OK",
	}, nil

}

func (s *Server) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	idStr := req.Book.Id

	book := s.Storage.GetOne(idStr)

	if !bookvalidator.IsValid(book) {
		return nil, fmt.Errorf("error: invalid data")
	}

	if book.Id != "" {
		book.Name = req.Book.Name
		book.Author = req.Book.Author
		book.Isbn = req.Book.Isbn
		book.Publisher = req.Book.Publisher
		book.Genre = req.Book.Genre
		book.YearOfPublication = req.Book.Year
		book.Pages = req.Book.Pages

		err := s.Storage.Update(book)

		if err != nil {
			return &pb.UpdateBookResponse{
				Result: "error",
			}, err
		}

		return &pb.UpdateBookResponse{
			Result: "OK",
		}, nil
	}

	return nil, fmt.Errorf("invalid id from request")
}
