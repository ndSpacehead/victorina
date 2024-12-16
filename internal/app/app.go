package app

import (
	"context"

	"victorina/internal/storage"
)

func Run() error {
	db, err := storage.New()
	if err != nil {
		return err
	}
	return db.Ping(context.Background())
}
