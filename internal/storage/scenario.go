package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"victorina/internal/model"
)

// CreateScenario writes new scenario to database.
func (repo *repository) CreateScenario(ctx context.Context, req model.CreateScenarioRequest) (uuid.UUID, error) {
	id := uuid.New()
	args := []any{
		id.String(),
		req.Name,
		req.Description,
	}
	query := `
INSERT INTO scenarios (external_id, "name", description)
VALUES (?, ?, ?);`
	_, err := repo.db.ExecContext(ctx, query, args...)
	return id, err
}

// ReadScenario returns scenario data from database by scenario's ID.
func (repo *repository) ReadScenario(ctx context.Context, id uuid.UUID) (*model.Scenario, error) {
	out := model.Scenario{
		ID: id,
	}
	query := `
SELECT "name", description
FROM scenarios
WHERE external_id = ?;`
	if err := repo.db.QueryRowContext(ctx, query, id.String()).Scan(&out.Name, &out.Description); err != nil {
		if errors.Is(errSQLite(err), ErrNoRows) {
			err = model.ErrScenarioNotFound
		}
		return nil, err
	}
	return &out, nil
}

// AllScenarios returns list of all stored scenarios.
func (repo *repository) AllScenarios(ctx context.Context) ([]model.Scenario, error) {
	query := `
SELECT external_id, "name", description
FROM scenarios;`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(errSQLite(err), ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var out []model.Scenario
	for rows.Next() {
		var (
			sc model.Scenario
			id string
		)
		if err := rows.Scan(&id, &sc.Name, &sc.Description); err != nil {
			return nil, err
		}
		sc.ID, err = uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, nil
}

// UpdateScenario updates question data.
func (repo *repository) UpdateScenario(ctx context.Context, req model.Scenario) error {
	args := []any{
		req.Name,
		req.Description,
		req.ID.String(),
	}
	query := `
UPDATE scenarios SET "name" = ?, description = ?
WHERE external_id = ?;`
	_, err := repo.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteScenario deletes scenario by ID.
func (repo *repository) DeleteScenario(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM scenarios WHERE external_id = ?;`
	_, err := repo.db.ExecContext(ctx, query, id.String())
	return err
}

func (repo *repository) ReadAssignedQuestion(ctx context.Context, req model.AssignedQuestionRequest) (*model.AssignedQuestion, error) {
	out := model.AssignedQuestion{
		Question: model.Question{
			ID: req.QuestionID,
		},
	}
	args := []any{
		req.ScenarioID.String(),
		req.QuestionID.String(),
	}
	query := `
SELECT q.question, q.answer, sq.score
FROM
	scenarios_questions sq
		INNER JOIN scenarios sc ON sc.external_id = ? AND sq.scenario_id = sc.id
		INNER JOIN questions q ON q.external_id = ? AND sq.question_id = q.id;`
	if err := repo.db.QueryRowContext(ctx, query, args...).Scan(&out.Q, &out.Answer, &out.Score); err != nil {
		if errors.Is(errSQLite(err), ErrNoRows) {
			err = model.ErrAssignedQuestionNotFound
		}
		return nil, err
	}
	return &out, nil
}

// AllAssignedQuestions returns list of all questions assigned to scenario.
func (repo *repository) AllAssignedQuestions(ctx context.Context, id uuid.UUID) ([]model.AssignedQuestion, error) {
	query := `
SELECT q.external_id, q.question, q.answer, sq.score
FROM
	scenarios_questions sq
		INNER JOIN scenarios sc ON sc.external_id = ? AND sq.scenario_id = sc.id
		INNER JOIN questions q ON sq.question_id = q.id;`
	rows, err := repo.db.QueryContext(ctx, query, id.String())
	if err != nil {
		if errors.Is(errSQLite(err), ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var out []model.AssignedQuestion
	for rows.Next() {
		var (
			aq model.AssignedQuestion
			id string
		)
		if err := rows.Scan(&id, &aq.Q, &aq.Answer, &aq.Score); err != nil {
			return nil, err
		}
		aq.ID, err = uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		out = append(out, aq)
	}
	return out, nil
}

// AllNotAssignedQuestions returns list of all questions not assigned to scenario.
func (repo *repository) AllNotAssignedQuestions(ctx context.Context, id uuid.UUID) ([]model.Question, error) {
	query := `
SELECT q.external_id, q.question, q.answer
FROM
	questions q LEFT JOIN (
		SELECT sq.question_id FROM scenarios_questions sq INNER JOIN scenarios sc ON sc.external_id = ? AND sq.scenario_id = sc.id
	)  t ON t.question_id = q.id
WHERE t.question_id IS NULL;`
	rows, err := repo.db.QueryContext(ctx, query, id.String())
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

// AssignQuestionToScenario assigns question to scenario.
func (repo *repository) AssignQuestionToScenario(ctx context.Context, req model.AssignQuestionRequest) error {
	args := []any{
		req.Score,
		req.ScenarioID.String(),
		req.QuestionID.String(),
	}
	query := `
INSERT INTO scenarios_questions (scenario_id, question_id, score)
SELECT sc.id, q.id, ? FROM scenarios sc, questions q WHERE sc.external_id = ? AND q.external_id = ?;`
	_, err := repo.db.ExecContext(ctx, query, args...)
	if errors.Is(errSQLite(err), ErrDuplicate) {
		return model.ErrAssignmentExists
	}
	return err
}

// UpdateAssignedQuestion updates assigned question data.
func (repo *repository) UpdateAssignedQuestion(ctx context.Context, req model.AssignQuestionRequest) error {
	args := []any{
		req.Score,
		req.ScenarioID.String(),
		req.QuestionID.String(),
	}
	query := `
UPDATE scenarios_questions SET score = ?
WHERE (scenario_id, question_id) IN (
	SELECT sc.id, q.id FROM scenarios sc, questions q WHERE sc.external_id = ? AND q.external_id = ?
);`
	_, err := repo.db.ExecContext(ctx, query, args...)
	return err
}

// ExcludeQuestionFromScenario excludes question from scenario.
func (repo *repository) ExcludeQuestionFromScenario(ctx context.Context, req model.ExcludeQuestionRequest) error {
	args := []any{
		req.ScenarioID.String(),
		req.QuestionID.String(),
	}
	query := `
DELETE FROM scenarios_questions
WHERE (scenario_id, question_id) IN (
	SELECT sc.id, q.id FROM scenarios sc, questions q WHERE sc.external_id = ? AND q.external_id = ?
);`
	_, err := repo.db.ExecContext(ctx, query, args...)
	return err
}
