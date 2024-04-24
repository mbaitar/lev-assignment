package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mbaitar/levenue-assignment/db"
	"log"
)

func createDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return db.CreateDatabase()
}

func main() {
	db, err := createDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tables := []string{
		"schema_migrations",
		"accounts",
		"integration_tokens",
		"subscriptions",
		"metrics",
		"trades",
	}

	for _, table := range tables {
		query := fmt.Sprintf("drop table if exists %s", table)
		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}
}
