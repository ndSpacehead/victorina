package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"victorina/internal/model"
)

func getScenarios(s *server, w http.ResponseWriter, r *http.Request) {
	scs, err := s.repo.AllScenarios(r.Context())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		writeInternalError(w, "Не удалось получить список сценариев")
		return
	}
	if err := s.tpl.render(w, "container", containerWithScenario(scs)); err != nil {
		writeRenderError(w)
	}
}

func getScenario(id uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	sc, err := s.repo.ReadScenario(r.Context(), id)
	if err != nil {
		writeNotFoundError(w, err, "Сценарий не найден", "Не удалось выполнить поиск сценария")
		return
	}
	if err := s.tpl.render(w, "scenario", *sc); err != nil {
		writeRenderError(w)
	}
}

func postScenario(s *server, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeBadRequest(w, "Не удалось прочитать данные формы")
		return
	}
	id, err := s.repo.CreateScenario(r.Context(), model.CreateScenarioRequest{
		Name:        r.PostFormValue("name"),
		Description: r.PostFormValue("description"),
	})
	if err != nil {
		writeInternalError(w, "Не удалось записать новый сценарий")
		return
	}
	sc, err := s.repo.ReadScenario(r.Context(), id)
	if err != nil {
		writeInternalError(w, "Не удалось прочитать новый сценарий")
		return
	}
	if err := s.tpl.render(w, "scenario-form", scenarioToSchema(model.Scenario{})); err != nil {
		writeRenderError(w)
		return
	}
	if err := s.tpl.render(w, "oob-scenario", scenarioToSchema(*sc)); err != nil {
		writeRenderError(w)
	}

}

func putScenario(id uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeBadRequest(w, "Не удалось прочитать данные формы")
		return
	}
	if err := s.repo.UpdateScenario(r.Context(), model.Scenario{
		ID:          id,
		Name:        r.PostFormValue("name"),
		Description: r.PostFormValue("description"),
	}); err != nil {
		writeInternalError(w, "Не удалось перезаписать сценарий")
		return
	}
	sc, err := s.repo.ReadScenario(r.Context(), id)
	if err != nil {
		writeNotFoundError(w, err, "Сценарий не найден", "Не удалось выполнить поиск сценария")
		return
	}
	if err := s.tpl.render(w, "scenario", scenarioToSchema(*sc)); err != nil {
		writeRenderError(w)
	}
}

func deleteScenario(id uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	if scid := r.Header.Get("X-Scenario-ID"); scid == id.String() {
		if sc, err := s.repo.ReadScenario(r.Context(), id); err == nil {
			sc.ID = uuid.Nil
			if err := s.tpl.render(w, "oob-scenario-form", scenarioToSchema(*sc)); err != nil {
				setHXTriggerHeader(w, dangerToast, "Не удалось обновить форму сценария")
			}
		}
		if err := s.tpl.render(w, "oob-question-part", scenariosQuestions{}); err != nil {
			setHXTriggerHeader(w, dangerToast, "Не удалось обновить форму сценария")
		}
	}
	if err := s.repo.DeleteScenario(r.Context(), id); err != nil {
		writeInternalError(w, "Не удалось удалить сценарий")
	}
}

func getScenariosQuestions(id uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
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
	}
}

func postScenariosQuestion(sid, qid uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeBadRequest(w, "Не удалось прочитать данные формы")
		return
	}
	score, err := strconv.Atoi(r.PostFormValue("score"))
	if err != nil {
		writeBadRequest(w, "Не удалось прочитать количество баллов")
		return
	}
	if err := s.repo.AssignQuestionToScenario(r.Context(), model.AssignQuestionRequest{
		ScenarioID: sid,
		QuestionID: qid,
		Score:      score,
	}); err != nil {
		switch {
		case errors.Is(err, model.ErrExists):
			writeBadRequest(w, "Сценарий уже содержит этот вопрос")
		default:
			writeInternalError(w, "Не удалось назначить вопрос сценарию")
		}
		return
	}
	aq, err := s.repo.ReadAssignedQuestion(r.Context(), model.AssignedQuestionRequest{
		ScenarioID: sid,
		QuestionID: qid,
	})
	if err != nil {
		writeInternalError(w, "Не удалось выполнить поиск вопроса сценария")
		return
	}
	if err := s.tpl.render(w, "oob-scenarios-question-form", scenarioQuestionSchema{}); err != nil {
		writeRenderError(w)
		return
	}
	if err := s.tpl.render(w, "oob-assigned-question", assignedQuestionToSchema(*aq, sid)); err != nil {
		writeRenderError(w)
	}
}

func putScenariosQuestion(sid, qid uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeBadRequest(w, "Не удалось прочитать данные формы")
		return
	}
	score, err := strconv.Atoi(r.PostFormValue("score"))
	if err != nil {
		writeBadRequest(w, "Не удалось прочитать количество баллов")
		return
	}
	if err := s.repo.UpdateAssignedQuestion(r.Context(), model.AssignQuestionRequest{
		ScenarioID: sid,
		QuestionID: qid,
		Score:      score,
	}); err != nil {
		writeInternalError(w, "Не удалось изменить данные назначенного сценарию вопроса")
		return
	}
	as, err := s.repo.ReadAssignedQuestion(r.Context(), model.AssignedQuestionRequest{
		ScenarioID: sid,
		QuestionID: qid,
	})
	if err != nil {
		writeNotFoundError(w, err, "Назначенный сценарию вопрос не найден", "Не удалось выполнить поиск назначенного сценарию вопроса")
		return
	}
	if err := s.tpl.render(w, "assigned-question", assignedQuestionToSchema(*as, sid)); err != nil {
		writeRenderError(w)
	}
}

func deleteScenariosQuestion(sid, qid uuid.UUID, s *server, w http.ResponseWriter, r *http.Request) {
	if id := r.Header.Get("X-Question-ID"); id == qid.String() {
		var sqs scenarioQuestionSchema
		if aq, err := s.repo.ReadAssignedQuestion(r.Context(), model.AssignedQuestionRequest{
			ScenarioID: sid,
			QuestionID: qid,
		}); err == nil {
			sqs.SID = sid.String()
			sqs.QID = qid.String()
			sqs.Question = aq.Q
			sqs.Score = aq.Score
		}
		if err := s.tpl.render(w, "oob-scenarios-question-form", sqs); err != nil {
			setHXTriggerHeader(w, dangerToast, "Не удалось обновить форму сценария")
		}
	}
	if err := s.repo.ExcludeQuestionFromScenario(r.Context(), model.ExcludeQuestionRequest{
		ScenarioID: sid,
		QuestionID: qid,
	}); err != nil {
		writeInternalError(w, "Не удалось исключить вопрос из сценария")
		return
	}
	q, err := s.repo.ReadQuestion(r.Context(), qid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return
		}
		writeInternalError(w, "Не удалось прочитать исключенный вопрос")
		return
	}
	if err := s.tpl.render(w, "oob-free-question", questionToSchema(*q)); err != nil {
		writeRenderError(w)
	}
}
