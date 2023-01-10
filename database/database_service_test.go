package database

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func dbConnection() (*sql.DB, error) {
	// Connect to the database
	db, err := ConnectToDB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestConnectToDB(t *testing.T) {
	db, err := dbConnection()
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
