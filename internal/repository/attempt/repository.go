package attempt

import (
	"context"

	"github.com/maximrakov/ai-quizzes-backend/internal/database/postgres"
	"github.com/maximrakov/ai-quizzes-backend/internal/model"
)

type repository struct {
	postgres *postgres.Postgres
}

func NewRepo(postgres *postgres.Postgres) *repository {
	return &repository{postgres: postgres}
}

func (r *repository) Create(ctx context.Context, attempt *model.Attempt) (*model.Attempt, error) {
	err := r.postgres.Pool.QueryRow(ctx,
		`INSERT INTO attempts (assignment_id, score, correct_question_ids, wrong_question_ids)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		attempt.AssignmentId,
		attempt.Score,
		attempt.CorrectQuestionIds,
		attempt.WrongQuestionIds,
	).Scan(&attempt.Id)
	if err != nil {
		return nil, err
	}
	return attempt, nil
}
