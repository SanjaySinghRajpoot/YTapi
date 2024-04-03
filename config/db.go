package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	// dsn := "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"
	connStr := "postgres://postgres:postgres@localhost/ytapi?sslmode=disable"
	// connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Connected to Database")

	DB = db

	return db, nil
}
