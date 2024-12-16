package storage

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB interface {
	Close(context.Context) error
	Ping(context.Context) error
}

type db struct {
	db *sql.DB
}

func New() (DB, error) {
	sqliteDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return &db{
		db: sqliteDB,
	}, nil
}

func (o *db) Close(_ context.Context) error {
	return o.db.Close()
}

func (o *db) Ping(ctx context.Context) error {
	return o.db.PingContext(ctx)
}
