package server

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"victorina/internal/model"
)

func getQuestions(s *server, w http.ResponseWriter, r *http.Request) {
	qs, err := s.repo.AllQuestions(r.Context())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		writeInternalError(w, "Не удалось получить список вопросов")
		return
	}
	if err := s.tpl.render(w, "container", containerWithQuestions(qs)); err != nil {
		writeRenderError(w)
	}
}

func getQuestion(id uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	q, err := s.repo.ReadQuestion(r.Context(), id)
	if err != nil {
		writeNotFoundError(w, err, "Вопрос не найден", "Не удалось выполнить поиск вопроса")
		return
	}
	if err := s.tpl.render(w, "question", *q); err != nil {
		writeRenderError(w)
	}
}

func postQuestion(s *server, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeBadRequest(w, "Не удалось прочитать данные формы")
		return
	}
	id, err := s.repo.CreateQuestion(r.Context(), model.CreateQuestionRequest{
		Q:      r.PostFormValue("question"),
		Answer: r.PostFormValue("answer"),
	})
	if err != nil {
		writeInternalError(w, "Не удалось записать новый вопрос")
		return
	}
	q, err := s.repo.ReadQuestion(r.Context(), id)
	if err != nil {
		writeInternalError(w, "Не удалось прочитать новый вопрос")
		return
	}
	if err := s.tpl.render(w, "question-form", questionToSchema(model.Question{})); err != nil {
		writeRenderError(w)
		return
	}
	if err := s.tpl.render(w, "oob-question", questionToSchema(*q)); err != nil {
		writeRenderError(w)
	}
}

func putQuestion(id uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeBadRequest(w, "Не удалось прочитать данные формы")
		return
	}
	if err := s.repo.UpdateQuestion(r.Context(), model.Question{
		ID:     id,
		Q:      r.PostFormValue("question"),
		Answer: r.PostFormValue("answer"),
	}); err != nil {
		writeInternalError(w, "Не удалось перезаписать вопрос")
		return
	}
	q, err := s.repo.ReadQuestion(r.Context(), id)
	if err != nil {
		writeNotFoundError(w, err, "Вопрос не найден", "Не удалось выполнить поиск вопроса")
		return
	}
	if err := s.tpl.render(w, "question", questionToSchema(*q)); err != nil {
		writeRenderError(w)
	}
}

func deleteQuestion(id uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	if qid := r.Header.Get("X-Question-ID"); qid == id.String() {
		if q, err := s.repo.ReadQuestion(r.Context(), id); err == nil {
			q.ID = uuid.Nil
			if err := s.tpl.render(w, "oob-question-form", questionToSchema(*q)); err != nil {
				setHXTriggerHeader(w, dangerToast, "Не удалось обновить форму вопроса")
			}
		}
	}
	if err := s.repo.DeleteQuestion(r.Context(), id); err != nil {
		writeInternalError(w, "Не удалось удалить вопрос")
	}
}
