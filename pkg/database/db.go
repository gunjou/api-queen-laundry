package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(dbUrl string) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	return dbpool
}