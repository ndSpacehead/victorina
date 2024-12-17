package server

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"strconv"

	"victorina/internal/storage"
)

//go:embed template
var indexHTML embed.FS

//go:embed static
var staticFiles embed.FS

// Server is a generic server contract.
type Server interface {
	// Serve starts listening and handling requests on incoming connections.
	Serve() error

	// Close gracefully shuts down the server.
	Close(context.Context) error
}

// Config - configuration values for HTTP-server.
type Config struct {
	Port uint16
	Repo storage.Repository
}

type server struct {
	repo storage.Repository
	srv  *http.Server
	tpl  *templates
}

// New is a constructor of HTTP-server.
func New(c Config) (Server, error) {
	if c.Repo == nil {
		return nil, errors.New("repository must be not nil")
	}
	tpl, err := newTemplates()
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	out := &server{
		repo: c.Repo,
		srv: &http.Server{
			Addr:    ":" + strconv.Itoa(int(c.Port)),
			Handler: mux,
		},
		tpl: tpl,
	}
	mux.Handle("/", newHomeHandler(out))
	mux.Handle("/ping", newPingHandler())
	mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))
	return out, nil
}

// Close gracefully shuts down the server.
func (s *server) Close(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// Serve starts listening and handling requests on incoming connections.
func (s *server) Serve() error {
	return s.srv.ListenAndServe()
}

func newPingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var code int
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			code = http.StatusOK
		default:
			code = http.StatusMethodNotAllowed
		}
		w.WriteHeader(code)
	}
}
