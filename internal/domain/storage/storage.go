package storage

import (
	"w-r-grpc/internal/domain/entity"
)

type Storage interface {
	GetOne(idStr string) entity.Book
	GetAll() []entity.Book
	Post(book entity.Book) error
	Update(book entity.Book) error
	Delete(idStr string) error
}
