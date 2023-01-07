package database

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func dbConnection() *sql.DB {
	// Connect to the database
	db, err := ConnectToDB()
	if err != nil {
		panic(err)
	}
	return db
}

func TestConnectToDB(t *testing.T) {
	assert.True(t, dbConnection() != nil)
}
