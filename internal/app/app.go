package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"victorina/internal/server"
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
	defer repo.Close(context.Background())
	if err := repo.Ping(context.Background()); err != nil {
		return err
	}
	srv, err := server.New(server.Config{
		Port: uint16(c.port),
		Repo: repo,
	})
	if err != nil {
		return err
	}
	sCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	g, gCtx := errgroup.WithContext(sCtx)
	g.Go(func() error {
		err := srv.Serve()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})
	g.Go(func() error {
		<-gCtx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := srv.Close(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				return err
			}
		}
		return nil
	})
	fmt.Printf("open UI: %s\n", srv.Address())
	return g.Wait()
}
