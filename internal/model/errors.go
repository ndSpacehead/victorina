package model

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is a common error when something not found.
	ErrNotFound = errors.New("not found")

	// ErrQuestionNotFound is an error for cases when a question not found.
	ErrQuestionNotFound = fmt.Errorf("question %w", ErrNotFound)

	// ErrAssignedQuestionNotFound is an error for cases when a question is not assigned to a scenario.
	ErrAssignedQuestionNotFound = fmt.Errorf("assigned %w", ErrQuestionNotFound)

	// ErrScenarioNotFound is an error for cases when a scenario not found.
	ErrScenarioNotFound = fmt.Errorf("scenario %w", ErrNotFound)

	// ErrExists is a common error when something already exists.
	ErrExists = errors.New("exists")

	// ErrAssignmentExists is an error for cases when a question assignment already existts.
	ErrAssignmentExists = fmt.Errorf("question assignment already %w", ErrExists)
)
