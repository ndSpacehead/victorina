package server

import (
	"github.com/google/uuid"

	"victorina/internal/model"
)

type scenarioSchema struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func scenarioToSchema(sc model.Scenario) scenarioSchema {
	var id string
	if sc.ID != uuid.Nil {
		id = sc.ID.String()
	}
	return scenarioSchema{
		ID:          id,
		Name:        sc.Name,
		Description: sc.Description,
	}
}

func scenariosToSchema(scs []model.Scenario) []scenarioSchema {
	out := make([]scenarioSchema, 0, len(scs))
	for _, sc := range scs {
		out = append(out, scenarioToSchema(sc))
	}
	return out
}
