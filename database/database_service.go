package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectToDB() (*sql.DB, error) {
	// Connect to the database
	db, err := sql.Open("postgres", "user=postgres password=1488322 dbname=goodsdb sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
