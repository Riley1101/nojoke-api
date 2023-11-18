package lib

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:admin@localhost:5432/nojoke?sslmode=disable")

	if err != nil {
		panic(err)
	}
	return db
}
