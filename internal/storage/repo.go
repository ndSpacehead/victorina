package storage

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Repository is an interface for DB.
type Repository interface {
	Close(context.Context) error
	Ping(context.Context) error
}

// Config for setup [Repository] implementation.
type Config struct {
	Filename string
}

// repository is a database handle providing operations for process logic.
type repository struct {
	db *sql.DB
}

func New(c Config) (Repository, error) {
	sqliteDB, err := sql.Open("sqlite3", c.Filename)
	if err != nil {
		return nil, err
	}
	return &repository{
		db: sqliteDB,
	}, nil
}

func (repo *repository) Close(_ context.Context) error {
	return repo.db.Close()
}

func (repo *repository) Ping(ctx context.Context) error {
	return repo.db.PingContext(ctx)
}
