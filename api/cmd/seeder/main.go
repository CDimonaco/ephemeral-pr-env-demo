package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type envConfig struct {
	DBConnection string
}

func mustGetEnvConfig() envConfig {
	db := os.Getenv("DB_CONNECTION")
	if db == "" {
		panic("DB_CONNECTION env is required")
	}
	return envConfig{DBConnection: db}
}

func main() {
	envConfig := mustGetEnvConfig()
	ctx := context.Background()
	seedPath := os.Args[1]

	log.Printf("seeder started - seed file: %s \n", seedPath)

	conn, err := pgx.Connect(ctx, envConfig.DBConnection)
	if err != nil {
		panic(err)
	}

	seedContent, err := os.ReadFile(seedPath)
	if err != nil {
		log.Panicf("error reading seed file: %s - %s", seedPath, err.Error())
	}

	_, err = conn.Exec(ctx, string(seedContent))
	if err != nil {
		log.Panicf("error executing seed file: %s - %s", seedPath, err.Error())
	}

	log.Printf("seeder stopped")
}
