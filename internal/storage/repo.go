package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"victorina/internal/model"
)

// Repository is an interface for DB.
type Repository interface {
	model.QuestionRepository

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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

func (repo *repository) init(ctx context.Context) error {
	return repo.upMigrations(ctx)
}

func (repo *repository) upMigrations(ctx context.Context) error {
	if _, err := repo.db.QueryContext(
		ctx,
		fmt.Sprintf(`SELECT name FROM sqlite_master WHERE type='table' AND name='%s';`, migrationTableName),
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if migrationErr := repo.upMigration(ctx, initialMigration); migrationErr != nil {
			return migrationErr
		}
	}
	var current int
	if err := repo.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT COALESCE(MAX(version), 0) FROM %s;`, migrationTableName)).Scan(&current); err != nil {
		return err
	}
	for i := current; i < len(migrations); i++ {
		if err := repo.upMigration(ctx, migrations[i]); err != nil {
			return err
		}
	}
	return nil
}

func (repo *repository) upMigration(ctx context.Context, query string) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, query); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}
