package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	"os"
)

var Bun *bun.DB

func CreateDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./levenue.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Init() error {
	db, err := CreateDatabase()
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	Bun = bun.NewDB(db, sqlitedialect.New())
	if len(os.Getenv("APP_DEBUG")) > 0 {
		Bun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	return nil
}
