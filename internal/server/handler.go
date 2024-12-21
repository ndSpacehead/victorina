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
			writeMethodNotAllowed(w, r.Method)
			return
		}
		if err := s.tpl.render(w, "index", nil); err != nil {
			writeRenderError(w)
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
			writeMethodNotAllowed(w, r.Method)
		}
	}
}

func newQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			writeNotFound(w, "Вопрос не найден")
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
			writeMethodNotAllowed(w, r.Method)
		}
	}
}

func newEditQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeMethodNotAllowed(w, r.Method)
			return
		}
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			writeNotFound(w, "Вопрос не найден")
			return
		}
		q, err := s.repo.ReadQuestion(r.Context(), id)
		if err != nil {
			writeNotFoundError(w, err, "Вопрос не найден", "Не удалось выполнить поиск вопроса")
			return
		}
		if err := s.tpl.render(w, "oob-question-form", questionToSchema(*q)); err != nil {
			writeRenderError(w)
		}
	}
}

func newNewQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeMethodNotAllowed(w, r.Method)
			return
		}
		if err := s.tpl.render(w, "oob-question-form", questionToSchema(model.Question{})); err != nil {
			writeRenderError(w)
		}
	}
}

func newGameHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeMethodNotAllowed(w, r.Method)
			return
		}
		questions, err := s.repo.AllQuestions(r.Context())
		if err != nil {
			writeInternalError(w, "Не удалось получить список вопросов")
			return
		}
		s.game.Reset(questions)
		if err := s.tpl.render(w, "container", containerWithGame(s.game.Scores())); err != nil {
			writeRenderError(w)
		}
	}
}

func newNextQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			writeMethodNotAllowed(w, r.Method)
			return
		}
		score, err := strconv.Atoi(r.PathValue("score"))
		if err != nil {
			writeNotFound(w, "Вопросов с оценкой %q нет", r.PathValue("score"))
			return
		}
		id, quantity, err := s.game.NextQuestion(score)
		if err != nil {
			writeNotFoundError(w, err, "Вопросы отсутствуют", "Не удалось получить вопрос")
			return
		}
		q, err := s.repo.ReadQuestion(r.Context(), id)
		if err != nil {
			writeNotFoundError(w, err, "Вопрос не найден", "Не удалось выполнить поиск вопроса")
			return
		}
		if err := s.tpl.render(w, "oob-current-question", questionToSchema(*q)); err != nil {
			writeRenderError(w)
			return
		}
		if quantity > 0 {
			if err := s.tpl.render(w, "card", strconv.Itoa(score)); err != nil {
				writeRenderError(w)
			}
		}
	}
}

func newPingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func writeBadRequest(w http.ResponseWriter, format string, args ...any) {
	setHXTriggerHeader(w, warningToast, format, args...)
	w.WriteHeader(http.StatusBadRequest)
}

func writeInternalError(w http.ResponseWriter, message string) {
	setHXTriggerHeader(w, dangerToast, message, nil...)
	w.WriteHeader(http.StatusInternalServerError)
}

func writeRenderError(w http.ResponseWriter) {
	writeInternalError(w, "Не удалось сформировать ответ")
}

func writeMethodNotAllowed(w http.ResponseWriter, method string) {
	setHXTriggerHeader(w, warningToast, "Метод обращения %q неприменим", method)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func writeNotFound(w http.ResponseWriter, format string, args ...any) {
	setHXTriggerHeader(w, infoToast, format, args...)
	w.WriteHeader(http.StatusNotFound)
}

func writeNotFoundError(w http.ResponseWriter, err error, main, alt string) {
	code := http.StatusInternalServerError
	constructor := dangerToast
	message := alt
	if errors.Is(err, model.ErrNotFound) {
		code = http.StatusNotFound
		constructor = warningToast
		message = main
	}
	setHXTriggerHeader(w, constructor, message, nil...)
	w.WriteHeader(code)
}
