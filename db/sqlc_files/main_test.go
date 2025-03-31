package db_sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	driver    = "postgres"
	sourseStr = "postgresql://admin:admin@localhost:60000/bank_db?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(driver, sourseStr)
	if err != nil {
		log.Fatalf("connection to DB: %v", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
