package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {
	var err error

	DB, err = pgxpool.New(
		context.Background(),
		"postgresql://neondb_owner:npg_aox0U6gNJKih@ep-odd-shape-al8hqpmx-pooler.c-3.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require",
	)

	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	log.Println("Connected to PostgreSQL")
}