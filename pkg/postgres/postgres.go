package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func OpenDB(user string, password string, host string, port string, dbname string) (*sql.DB, error) {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	db, err := sql.Open("postgres", connStr)

	return db, err
}
