package server

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"victorina/internal/model"
)

func newHomeHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if err := s.tpl.render(w, "index", nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func newQuestionsHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getQuestions(s, w, r)
		case http.MethodPost:
			postQuestion(s, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func newQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		switch r.Method {
		case http.MethodGet:
			getQuestion(id, s, w, r)
		case http.MethodPut:
			putQuestion(id, s, w, r)
		case http.MethodDelete:
			deleteQuestion(id, s, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func newEditQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		q, err := s.repo.ReadQuestion(r.Context(), id)
		if err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, model.ErrNotFound) {
				code = http.StatusNotFound
			}
			w.WriteHeader(code)
			return
		}
		if err := s.tpl.render(w, "oob-question-form", questionToSchema(*q)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func newNewQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := s.tpl.render(w, "oob-question-form", questionToSchema(model.Question{})); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
