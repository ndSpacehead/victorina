package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"victorina/internal/model"
)

// CreateQuestion writes new question to database.
func (repo *repository) CreateQuestion(ctx context.Context, req model.CreateQuestionRequest) (uuid.UUID, error) {
	id := uuid.New()
	args := []any{
		id.String(),
		req.Q,
		req.Answer,
	}
	query := `
INSERT INTO questions (external_id, question, answer)
VALUES (?, ?, ?);`
	_, err := repo.db.ExecContext(ctx, query, args...)
	return id, err
}

// ReadQuestion returns question data from database by question's ID.
func (repo *repository) ReadQuestion(ctx context.Context, id uuid.UUID) (*model.Question, error) {
	out := model.Question{
		ID: id,
	}
	query := `
SELECT question, answer
FROM questions
WHERE external_id = ?;`
	if err := repo.db.QueryRowContext(ctx, query, id.String()).Scan(&out.Q, &out.Answer); err != nil {
		if errors.Is(errSQLite(err), ErrNoRows) {
			err = model.ErrQuestionNotFound
		}
		return nil, err
	}
	return &out, nil
}

// AllQuestions returns list of all stored questions.
func (repo *repository) AllQuestions(ctx context.Context) ([]model.Question, error) {
	query := `
SELECT external_id, question, answer
FROM questions;`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(errSQLite(err), ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var out []model.Question
	for rows.Next() {
		var (
			q  model.Question
			id string
		)
		if err := rows.Scan(&id, &q.Q, &q.Answer); err != nil {
			return nil, err
		}
		q.ID, err = uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, nil
}

// UpdateQuestion updates question data.
func (repo *repository) UpdateQuestion(ctx context.Context, req model.Question) error {
	args := []any{
		req.Q,
		req.Answer,
		req.ID.String(),
	}
	query := `
UPDATE questions SET question = ?, answer = ?
WHERE external_id = ?;`
	_, err := repo.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteQuestion deletes question by ID.
func (repo *repository) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM questions WHERE external_id = ?;`
	_, err := repo.db.ExecContext(ctx, query, id.String())
	return err
}
