package server

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"victorina/internal/model"
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

	// Address returns TCP address for the server to listen on.
	Address() string
}

// Config - configuration values for HTTP-server.
type Config struct {
	// Port specifies the TCP port for the server to listen on.
	Port uint16

	// Repo is a DB repository.
	Repo storage.Repository
}

type server struct {
	repo storage.Repository
	srv  *http.Server
	tpl  *templates
	game *model.Game
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
		tpl:  tpl,
		game: new(model.Game),
	}
	mux.Handle("/", newHomeHandler(out))
	mux.Handle("/ping", newPingHandler())
	mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))
	mux.Handle("/questions", newQuestionsHandler(out))
	mux.Handle("/questions/{id}", newQuestionHandler(out))
	mux.Handle("/questions/{id}/edit", newEditQuestionHandler(out))
	mux.Handle("/questions/{id}/answer", newAnswerHandler(out))
	mux.Handle("/questions/new", newNewQuestionHandler(out))
	mux.Handle("/questions/{id}/assign", newAssignQuestionHandler(out))
	mux.Handle("/scenarios", newScenariosHandler(out))
	mux.Handle("/scenarios/{id}", newScenarioHandler(out))
	mux.Handle("/scenarios/{id}/edit", newEditScenarioHandler(out))
	mux.Handle("/scenarios/{id}/questions", newScenariosQuestionsHandler(out))
	mux.Handle("/scenarios/{sid}/questions/{qid}", newScenariosQuestionHandler(out))
	mux.Handle("/scenarios/{sid}/questions/{qid}/edit", newEditScenariosQuestionHandler(out))
	mux.Handle("/scenarios/new", newNewScenarioHandler(out))
	mux.Handle("/game", newGameHandler(out))
	mux.Handle("/game/scenarios/{id}", newScenarioGameHandler(out))
	mux.Handle("/game/next/{score}", newNextQuestionHandler(out))
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

// Address returns TCP address for the server to listen on.
func (s *server) Address() string {
	local := "http://127.0.0.1"
	if s.srv.Addr == "" {
		return local
	}
	parts := strings.Split(s.srv.Addr, ":")
	if len(parts) < 2 {
		return local
	}
	if parts[0] == "" {
		parts[0] = local
	}
	return strings.Join(parts, ":")
}
