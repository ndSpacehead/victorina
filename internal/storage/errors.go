package storage

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate = errors.New("duplicate constraint")
	ErrNoRows    = errors.New("no rows in result set")
)

func errSQLite(err error) error {
	var sqliteErr sqlite3.Error
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNoRows
	case errors.As(err, &sqliteErr):
		if errors.Is(err, sqlite3.ErrConstraintUnique) || errors.Is(err, sqlite3.ErrConstraintPrimaryKey) {
			return ErrDuplicate
		}
		fallthrough
	default:
		return err
	}
}
