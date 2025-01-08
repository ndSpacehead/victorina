package model

import (
	"context"

	"github.com/google/uuid"
)

// ScenarioRepository is a CRUD API for scenario entity.
type ScenarioRepository interface {
	// CreateScenario writes new scenario's data.
	CreateScenario(context.Context, CreateScenarioRequest) (uuid.UUID, error)

	// ReadScenario returns scenario data by scenario's ID.
	ReadScenario(context.Context, uuid.UUID) (*Scenario, error)

	// AllScenarios returns list of all stored scenarios.
	AllScenarios(context.Context) ([]Scenario, error)

	// UpdateScenario updates question data.
	UpdateScenario(context.Context, Scenario) error

	// DeleteScenario deletes scenario by ID.
	DeleteScenario(context.Context, uuid.UUID) error

	// ReadAssignedQuestion returns the question assigned to the scenario.
	ReadAssignedQuestion(context.Context, AssignedQuestionRequest) (*AssignedQuestion, error)

	// AllAssignedQuestions returns list of all questions assigned to scenario.
	AllAssignedQuestions(context.Context, uuid.UUID) ([]AssignedQuestion, error)

	// AllNotAssignedQuestions returns list of all questions not assigned to scenario.
	AllNotAssignedQuestions(context.Context, uuid.UUID) ([]Question, error)

	// AssignQuestionToScenario assigns question to scenario.
	AssignQuestionToScenario(context.Context, AssignQuestionRequest) error

	// UpdateAssignedQuestion updates assigned question data.
	UpdateAssignedQuestion(context.Context, AssignQuestionRequest) error

	// ExcludeQuestionFromScenario excludes question from scenario.
	ExcludeQuestionFromScenario(context.Context, ExcludeQuestionRequest) error
}

// Scenario is a scenario entity.
type Scenario struct {
	// ID is a scenario's ID.
	ID uuid.UUID

	// Name is a scenario's title.
	Name string

	// Description is a scenario's description.
	Description string
}

// CreateScenarioRequest is an object with scenario data for creating scenario entity.
type CreateScenarioRequest struct {
	// Name is a scenario's title.
	Name string

	// Description is a scenario's description.
	Description string
}

// AssignedQuestion is a question entity which assigned to the scenario.
type AssignedQuestion struct {
	// Question is a question entity.
	Question

	// Score is a question grade in the scenario.
	Score int
}

// AssignQuestionRequest is an object for assign the question to the scenario.
type AssignQuestionRequest struct {
	// ScenarioID is a scenario's ID.
	ScenarioID uuid.UUID

	// QuestionID is a question's ID.
	QuestionID uuid.UUID

	// Score is a question's grade in the scenario.
	Score int
}

// ExcludeQuestionRequest is an object for exclude a question from a scenario.
type ExcludeQuestionRequest struct {
	// ScenarioID is a scenario's ID.
	ScenarioID uuid.UUID

	// QuestionID is a question's ID.
	QuestionID uuid.UUID
}

// AssignedQuestionRequest is an object for read an assigned question to a scenario.
type AssignedQuestionRequest struct {
	// ScenarioID is a scenario's ID.
	ScenarioID uuid.UUID

	// QuestionID is a question's ID.
	QuestionID uuid.UUID
}
