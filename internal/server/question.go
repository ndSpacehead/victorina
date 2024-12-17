package server

import (
	"errors"
	"net/http"
	"strconv"

	"victorina/internal/model"

	"github.com/google/uuid"
)

func getQuestions(s *server, w http.ResponseWriter, r *http.Request) {
	qs, err := s.repo.AllQuestions(r.Context())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := s.tpl.render(w, "container", containerWithQuestions(qs)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getQuestion(id uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	q, err := s.repo.ReadQuestion(r.Context(), id)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, model.ErrNotFound) {
			code = http.StatusNotFound
		}
		w.WriteHeader(code)
		return
	}
	if err := s.tpl.render(w, "question", *q); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func postQuestion(s *server, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	score, err := strconv.Atoi(r.PostFormValue("score"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := s.repo.CreateQuestion(r.Context(), model.CreateQuestionRequest{
		Q:      r.PostFormValue("question"),
		Answer: r.PostFormValue("answer"),
		Score:  score,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	q, err := s.repo.ReadQuestion(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := s.tpl.render(w, "question-form", questionToSchema(model.Question{})); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := s.tpl.render(w, "oob-question", questionToSchema(*q)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func putQuestion(id uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	score, err := strconv.Atoi(r.PostFormValue("score"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := s.repo.UpdateQuestion(r.Context(), model.Question{
		ID:     id,
		Q:      r.PostFormValue("question"),
		Answer: r.PostFormValue("answer"),
		Score:  score,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	if err := s.tpl.render(w, "question", questionToSchema(*q)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteQuestion(id uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	if err := s.repo.DeleteQuestion(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
