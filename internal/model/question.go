package model

import (
	"context"

	"github.com/google/uuid"
)

// QuestionRepository is a CRUD API for question entity.
type QuestionRepository interface {
	// CreateQuestion writes new question to database.
	CreateQuestion(context.Context, CreateQuestionRequest) (uuid.UUID, error)

	// ReadQuestion returns question data from database by question ID.
	ReadQuestion(context.Context, uuid.UUID) (*Question, error)

	// AllQuestions returns list of all stored questions.
	AllQuestions(context.Context) ([]Question, error)

	// UpdateQuestion updates question data.
	UpdateQuestion(context.Context, Question) error

	// DeleteQuestion deletes question by ID.
	DeleteQuestion(context.Context, uuid.UUID) error
}

// Question is a question entity.
type Question struct {
	ID     uuid.UUID
	Q      string
	Answer string
}

// CreateQuestionRequest is an object with question data for creating question entity.
type CreateQuestionRequest struct {
	Q      string
	Answer string
}
