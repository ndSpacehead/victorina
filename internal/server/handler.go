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
		scs, err := s.repo.AllScenarios(r.Context())
		if err != nil {
			writeInternalError(w, "Не удалось получить список сценариев")
			return
		}
		if err := s.tpl.render(w, "index", scenariosToSchema(scs)); err != nil {
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

func newAnswerHandler(s *server) http.HandlerFunc {
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
		if err := s.tpl.render(w, "showed-answer", q.Answer); err != nil {
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

func newAssignQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeMethodNotAllowed(w, r.Method)
			return
		}
		qid, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			writeNotFound(w, "Вопрос не найден")
			return
		}
		sid, err := uuid.Parse(r.Header.Get("X-Scenario-ID"))
		if err != nil {
			writeNotFound(w, "Сценарий не найден")
			return
		}
		sc, err := s.repo.ReadScenario(r.Context(), sid)
		if err != nil {
			writeNotFoundError(w, err, "Сценарий не найден", "Не удалось выполнить поиск сценария.")
			return
		}
		q, err := s.repo.ReadQuestion(r.Context(), qid)
		if err != nil {
			writeNotFoundError(w, err, "Вопрос не найден", "Не удалось выполнить поиск вопроса.")
			return
		}
		qs := questionToSchema(*q)
		if err := s.tpl.render(w, "oob-scenarios-question-form", scenarioQuestionSchema{
			SID:      sc.ID.String(),
			QID:      q.ID.String(),
			Question: qs.Question,
		}); err != nil {
			writeRenderError(w)
		}
	}
}

func newScenariosHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getScenarios(s, w, r)
		case http.MethodPost:
			postScenario(s, w, r)
		default:
			writeMethodNotAllowed(w, r.Method)
		}
	}
}

func newScenarioHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			writeNotFound(w, "Сценарий не найден")
			return
		}
		switch r.Method {
		case http.MethodGet:
			getScenario(id, s, w, r)
		case http.MethodPut:
			putScenario(id, s, w, r)
		case http.MethodDelete:
			deleteScenario(id, s, w, r)
		default:
			writeMethodNotAllowed(w, r.Method)
		}
	}
}

func newEditScenarioHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeMethodNotAllowed(w, r.Method)
			return
		}
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			writeNotFound(w, "Сценарий не найден")
			return
		}
		sc, err := s.repo.ReadScenario(r.Context(), id)
		if err != nil {
			writeNotFoundError(w, err, "Сценарий не найден", "Не удалось выполнить поиск сценария")
			return
		}
		aqs, err := s.repo.AllAssignedQuestions(r.Context(), id)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			writeInternalError(w, "Не удалось получить список вопросов для сценария")
			return
		}
		nqs, err := s.repo.AllNotAssignedQuestions(r.Context(), id)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			writeInternalError(w, "Не удалось получить список вопросов для сценария")
			return
		}
		if err := s.tpl.render(w, "oob-question-part", scenariosQuestions{
			AssignedList: assignedQuestionsToSchema(aqs, id),
			FreeList:     questionsToSchema(nqs),
		}); err != nil {
			writeRenderError(w)
			return
		}
		if err := s.tpl.render(w, "oob-scenario-form", scenarioToSchema(*sc)); err != nil {
			writeRenderError(w)
		}
	}
}

func newNewScenarioHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeMethodNotAllowed(w, r.Method)
			return
		}
		if err := s.tpl.render(w, "oob-question-part", scenariosQuestions{}); err != nil {
			writeRenderError(w)
			return
		}
		if err := s.tpl.render(w, "oob-scenario-form", scenarioToSchema(model.Scenario{})); err != nil {
			writeRenderError(w)
		}
	}
}

func newScenariosQuestionsHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			writeNotFound(w, "Сценарий не найден")
			return
		}
		switch r.Method {
		case http.MethodGet:
			getScenariosQuestions(id, s, w, r)
		default:
			writeMethodNotAllowed(w, r.Method)
		}
	}
}

func newScenariosQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid, err := uuid.Parse(r.PathValue("sid"))
		if err != nil {
			writeNotFound(w, "Сценарий не найден")
			return
		}
		qid, err := uuid.Parse(r.PathValue("qid"))
		if err != nil {
			writeNotFound(w, "Вопрос не найден")
			return
		}
		switch r.Method {
		case http.MethodPost:
			postScenariosQuestion(sid, qid, s, w, r)
		case http.MethodDelete:
			deleteScenariosQuestion(sid, qid, s, w, r)
		case http.MethodPut:
			putScenariosQuestion(sid, qid, s, w, r)
		default:
			writeMethodNotAllowed(w, r.Method)
		}
	}
}

func newEditScenariosQuestionHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid, err := uuid.Parse(r.PathValue("sid"))
		if err != nil {
			writeNotFound(w, "Сценарий не найден")
			return
		}
		qid, err := uuid.Parse(r.PathValue("qid"))
		if err != nil {
			writeNotFound(w, "Вопрос не найден")
			return
		}
		aq, err := s.repo.ReadAssignedQuestion(r.Context(), model.AssignedQuestionRequest{
			ScenarioID: sid,
			QuestionID: qid,
		})
		if err != nil {
			writeNotFoundError(w, err, "Вопрос не назначен сценарию", "Не удалось прочитать назначенный сценарию вопрос")
			return
		}
		if err := s.tpl.render(w, "oob-scenarios-question-form", scenarioQuestionSchema{
			SID:      sid.String(),
			QID:      qid.String(),
			Question: aq.Q,
			Score:    aq.Score,
			Assigned: true,
		}); err != nil {
			writeRenderError(w)
		}
	}
}

func newGameHandler(_ *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHXTriggerHeader(w, infoToast, "Устаревший режим игры не доступен")
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func newScenarioGameHandler(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeMethodNotAllowed(w, r.Method)
			return
		}
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			writeNotFound(w, "Сценарий не найден")
			return
		}
		sc, err := s.repo.ReadScenario(r.Context(), id)
		if err != nil {
			writeNotFoundError(w, err, "Сценарий не найден", "Не удалось прочитать сценарий")
			return
		}
		questions, err := s.repo.AllAssignedQuestions(r.Context(), id)
		if err != nil {
			writeInternalError(w, "Не удалось получить список вопросов")
			return
		}
		s.game.Reset(questions)
		if err := s.tpl.render(w, "container", containerWithGame(sc.Name, s.game.Scores())); err != nil {
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
		if err := s.tpl.render(w, "oob-current-question", questionToGameQuestionSchema(*q, score)); err != nil {
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
