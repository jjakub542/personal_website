package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var TestDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5432", "jjakub542", "123", "db_test")
	TestDB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	code := m.Run()
	TestDB.Close()
	os.Exit(code)
}
