package dblogic

import (
	"database/sql"
	"log"
	"w-r-grpc/internal/domain/entity"
)

import (
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}

func (d *Database) GetOne(idStr string) entity.Book {

	rows, err := d.db.Query("SELECT * FROM books WHERE id=$1", idStr)
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

func (d *Database) GetAll() []entity.Book {

	rows, err := d.db.Query("SELECT * FROM books")
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

func (d *Database) Post(book entity.Book) error {
	_, err := d.db.Exec(`INSERT INTO books (id, name, author, isbn, publisher, genre, year_of_publication, pages) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		book.Id, book.Name, book.Author, book.Isbn, book.Publisher, book.Genre, book.YearOfPublication, book.Pages)

	return err
}

func (d *Database) Update(book entity.Book) error {
	_, err := d.db.Exec("UPDATE books SET name=$2, author=$3, isbn=$4, publisher=$5, genre=$6, year_of_publication=$7, pages=$8 WHERE id=$1",
		book.Id, book.Name, book.Author, book.Isbn, book.Publisher, book.Genre, book.YearOfPublication, book.Pages)

	return err

}

func (d *Database) Delete(idStr string) error {
	_, err := d.db.Exec("DELETE FROM books WHERE id=$1", idStr)

	return err
}
