package server

import (
	"github.com/google/uuid"

	"victorina/internal/model"
)

type questionSchema struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func questionToSchema(q model.Question) questionSchema {
	var id string
	if q.ID != uuid.Nil {
		id = q.ID.String()
	}
	return questionSchema{
		ID:       id,
		Question: q.Q,
		Answer:   q.Answer,
	}
}

func questionsToSchema(qs []model.Question) []questionSchema {
	out := make([]questionSchema, 0, len(qs))
	for _, q := range qs {
		out = append(out, questionToSchema(q))
	}
	return out
}

type scenarioQuestionSchema struct {
	SID      string `json:"sid"`
	QID      string `json:"qid"`
	Question string `json:"question"`
	Score    int    `json:"score"`
	Assigned bool   `json:"assigned"`
}

type assignedQuestionSchema struct {
	SID      string `json:"sid"`
	QID      string `json:"qid"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Score    int    `json:"score"`
}

func assignedQuestionToSchema(aq model.AssignedQuestion, sid uuid.UUID) assignedQuestionSchema {
	var qid string
	if aq.ID != uuid.Nil {
		qid = aq.ID.String()
	}
	return assignedQuestionSchema{
		SID:      sid.String(),
		QID:      qid,
		Question: aq.Q,
		Answer:   aq.Answer,
		Score:    aq.Score,
	}
}

func assignedQuestionsToSchema(aqs []model.AssignedQuestion, sid uuid.UUID) []assignedQuestionSchema {
	out := make([]assignedQuestionSchema, 0, len(aqs))
	for _, aq := range aqs {
		out = append(out, assignedQuestionToSchema(aq, sid))
	}
	return out
}
