package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func DBstart() {
	cfg := LoadDBConfig()
	log.Println("Connecting with:", cfg.connString())

	poolConfig, err := pgxpool.ParseConfig(cfg.connString())
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection %v", err)
	}

	fmt.Println("Connected to Postgres")

	var currentTime time.Time
	DB.QueryRow(ctx, "SELECT NOW()").Scan(&currentTime)
	fmt.Println(currentTime)
}
