package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Repository is an interface for DB.
type Repository interface {
	// Close closes the database and prevents new queries from starting.
	Close(context.Context) error

	// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
	Ping(context.Context) error
}

// Config for setup [Repository] implementation.
type Config struct {
	Filename string
}

// repository is a database handle providing operations for CRUD.
type repository struct {
	db *sql.DB
}

// New is a constructor of [Repository] implementation with SQLite.
func New(c Config) (Repository, error) {
	db, err := sql.Open("sqlite3", c.Filename)
	if err != nil {
		return nil, err
	}
	out := &repository{
		db: db,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := out.init(ctx); err != nil {
		out.Close(context.Background())
		return nil, err
	}
	return out, nil
}

// Close closes the database and prevents new queries from starting.
func (repo *repository) Close(_ context.Context) error {
	return repo.db.Close()
}

// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
func (repo *repository) Ping(ctx context.Context) error {
	return repo.db.PingContext(ctx)
}

// init is naive implementation of migrations mechanism.
func (repo *repository) init(ctx context.Context) error {
	query := `
CREATE TABLE IF NOT EXISTS questions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	external_id VARCHAR(36) NOT NULL,
	question TEXT NOT NULL,
	answer TEXT NOT NULL,
	score INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS questions_external_id_idx ON questions (external_id);`
	_, err := repo.db.ExecContext(ctx, query)
	return err
}
