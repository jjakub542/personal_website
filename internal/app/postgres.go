package app

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

var (
	dbname = os.Getenv("DB_NAME")
	// dbnametest         = os.Getenv("DB_NAME_TEST")
	password       = os.Getenv("DB_PASSWORD")
	username       = os.Getenv("DB_USERNAME")
	port           = os.Getenv("DB_PORT")
	host           = os.Getenv("DB_HOST")
	postgresClient *pgxpool.Pool
)

func PostgresClient() *pgxpool.Pool {
	var err error
	if postgresClient != nil {
		return postgresClient
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)
	db, errConnecting := pgxpool.New(context.Background(), connStr)
	if errConnecting != nil {
		log.Fatalf("Postgres connection failed: %v", err)
	}
	if err = db.Ping(context.Background()); err != nil {
		log.Fatalf("Postgres ping failed: %v", err)
	}
	fmt.Println("✅ Connected to Postgres")
	return db
}
