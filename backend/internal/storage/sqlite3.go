package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func BootstrapSqlite3(filePath string, timeout time.Duration) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(timeout)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func CloseSqlite3(db *sql.DB) error {
	return db.Close()
}
