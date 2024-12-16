package app

import (
	"context"

	"victorina/internal/storage"
)

func Run() error {
	c := readConfig()
	repo, err := storage.New(storage.Config{
		Filename: c.dbFilename,
	})
	if err != nil {
		return err
	}
	return repo.Ping(context.Background())
}
