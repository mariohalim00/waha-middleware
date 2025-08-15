package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context) *pgxpool.Pool {

	var (
		database = os.Getenv("DB_DATABASE")
		password = os.Getenv("DB_PASSWORD")
		username = os.Getenv("DB_USERNAME")
		port     = os.Getenv("DB_PORT")
		host     = os.Getenv("DB_HOST")
		schema   = os.Getenv("DB_SCHEMA")
	)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	return pool
}
