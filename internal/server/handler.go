package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"victorina/internal/model"
)

func newHomeHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := s.tpl.render(w, "index", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

func newGameHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		questions, err := s.repo.AllQuestions(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.game.Reset(questions)
		if err := s.tpl.render(w, "container", containerWithGame(s.game.Scores())); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func newNextQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		score, err := strconv.Atoi(r.PathValue("score"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id, quantity, err := s.game.NextQuestion(score)
		if err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, model.ErrNotFound) {
				code = http.StatusNotFound
			}
			w.WriteHeader(code)
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
		if err := s.tpl.render(w, "oob-current-question", questionToSchema(*q)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if quantity > 0 {
			if err := s.tpl.render(w, "card", strconv.Itoa(score)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
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
