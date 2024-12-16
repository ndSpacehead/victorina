package app

import (
	"context"

	"victorina/internal/storage"
)

// Run is a point of initialize, start and shutdown application components.
func Run() error {
	c, err := readConfig()
	if err != nil {
		return err
	}
	repo, err := storage.New(storage.Config{
		Filename: c.dbFilename,
	})
	if err != nil {
		return err
	}
	return repo.Ping(context.Background())
}
