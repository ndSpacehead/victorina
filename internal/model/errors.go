package model

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is a common error when something not found.
	ErrNotFound = errors.New("not found")

	// ErrQuestionNotFound is an error for cases when question not found.
	ErrQuestionNotFound = fmt.Errorf("question %w", ErrNotFound)
)
