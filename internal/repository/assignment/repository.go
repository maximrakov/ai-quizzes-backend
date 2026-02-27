package assignment

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

func (r *repository) Create(ctx context.Context, assignment *model.Assignment) (*model.Assignment, error) {
	err := r.postgres.Pool.QueryRow(ctx,
		"INSERT INTO assignments (quiz_id, student_id) VALUES ($1, $2) RETURNING id",
		assignment.QuizId, assignment.StudentId,
	).Scan(&assignment.Id)
	if err != nil {
		return nil, err
	}
	return assignment, nil
}

func (r *repository) FindById(ctx context.Context, id int) (*model.Assignment, error) {
	a := &model.Assignment{}
	err := r.postgres.Pool.QueryRow(ctx,
		"SELECT id, quiz_id, student_id FROM assignments WHERE id = $1", id,
	).Scan(&a.Id, &a.QuizId, &a.StudentId)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *repository) FindByStudentId(ctx context.Context, studentId int) ([]*model.Assignment, error) {
	rows, err := r.postgres.Pool.Query(ctx,
		"SELECT id, quiz_id, student_id FROM assignments WHERE student_id = $1", studentId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*model.Assignment
	for rows.Next() {
		a := &model.Assignment{}
		if err = rows.Scan(&a.Id, &a.QuizId, &a.StudentId); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, nil
}
