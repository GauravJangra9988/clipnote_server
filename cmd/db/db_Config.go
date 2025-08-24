package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	User            string
	Password        string
	Host            string
	Port            string
	Name            string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
}

func LoadDBConfig() DBConfig {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	maxConns := int32(10)
	minConns := int32(2)
	maxConnLifetime := 30 * time.Minute

	return DBConfig{
		User:            os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		Host:            os.Getenv("DB_HOST"),
		Port:            os.Getenv("DB_PORT"),
		Name:            os.Getenv("DB_NAME"),
		SSLMode:         os.Getenv("DB_SSLMODE"),
		MaxConns:        maxConns,
		MinConns:        minConns,
		MaxConnLifetime: maxConnLifetime,
	}
}

func (cfg DBConfig) connString() string {

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
}
