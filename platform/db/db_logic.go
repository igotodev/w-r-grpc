package db

import (
	"database/sql"
	"fmt"
	"log"
	"w-r-grpc/config"
	"w-r-grpc/platform/entity"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func OpenDB() {
	cfg := config.InitConfig()
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", cfg.Driver, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err = sql.Open(cfg.Driver, connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	if err = db.Close(); err != nil {
		log.Fatal(err)
	}
}

func Select(idStr string) entity.Book {

	rows, err := db.Query("SELECT * FROM books WHERE id=$1", idStr)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	var book entity.Book

	for rows.Next() {
		err = rows.Scan(&book.Id, &book.Name, &book.Author, &book.Isbn, &book.Publisher, &book.Genre, &book.YearOfPublication, &book.Pages)
		if err != nil {
			log.Println(err)
		}
	}

	return book
}

func SelectAll() []entity.Book {

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	var books []entity.Book

	for rows.Next() {
		var book entity.Book

		err = rows.Scan(&book.Id, &book.Name, &book.Author, &book.Isbn, &book.Publisher, &book.Genre, &book.YearOfPublication, &book.Pages)
		if err != nil {
			log.Println(err)
		}

		books = append(books, book)
	}

	return books
}

func Insert(book entity.Book) {
	_, err = db.Exec(`INSERT INTO books (id, name, author, isbn, publisher, genre, year_of_publication, pages) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		book.Id, book.Name, book.Author, book.Isbn, book.Publisher, book.Genre, book.YearOfPublication, book.Pages)

	if err != nil {
		log.Println(err)
	}
}

func Update(book entity.Book) {
	_, err = db.Exec("UPDATE books SET name=$2, author=$3, isbn=$4, publisher=$5, genre=$6, year_of_publication=$7, pages=$8 WHERE id=$1",
		book.Id, book.Name, book.Author, book.Isbn, book.Publisher, book.Genre, book.YearOfPublication, book.Pages)

	if err != nil {
		log.Println(err)
	}

}

func Delete(idStr string) bool {
	_, err := db.Exec("DELETE FROM books WHERE id=$1", idStr)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
