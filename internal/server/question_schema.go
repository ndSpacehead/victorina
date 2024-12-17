package server

import (
	"github.com/google/uuid"

	"victorina/internal/model"
)

type questionSchema struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Score    int    `json:"score"`
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
		Score:    q.Score,
	}
}

func questionsToSchema(qs []model.Question) []questionSchema {
	out := make([]questionSchema, 0, len(qs))
	for _, q := range qs {
		out = append(out, questionToSchema(q))
	}
	return out
}
