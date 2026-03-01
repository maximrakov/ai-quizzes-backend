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

func (r *repository) FindByQuizId(ctx context.Context, quizId int) ([]*model.Attempt, error) {
	rows, err := r.postgres.Pool.Query(ctx,
		`SELECT a.id, a.assignment_id, a.score, a.correct_question_ids, a.wrong_question_ids,
		        u.id, u.nickname,
		        q.id, q.title
		 FROM attempts a
		 JOIN assignments asgn ON a.assignment_id = asgn.id
		 JOIN users u ON asgn.student_id = u.id
		 JOIN quizzes q ON asgn.quiz_id = q.id
		 WHERE asgn.quiz_id = $1`,
		quizId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*model.Attempt
	for rows.Next() {
		a := &model.Attempt{User: &model.User{}, Quiz: &model.Quiz{}}
		if err = rows.Scan(
			&a.Id, &a.AssignmentId, &a.Score, &a.CorrectQuestionIds, &a.WrongQuestionIds,
			&a.User.Id, &a.User.Nickname,
			&a.Quiz.Id, &a.Quiz.Title,
		); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, nil
}

func (r *repository) FindByStudentId(ctx context.Context, studentId int) ([]*model.Attempt, error) {
	rows, err := r.postgres.Pool.Query(ctx,
		`SELECT a.id, a.assignment_id, a.score, a.correct_question_ids, a.wrong_question_ids,
		        u.id, u.nickname,
		        q.id, q.title
		 FROM attempts a
		 JOIN assignments asgn ON a.assignment_id = asgn.id
		 JOIN users u ON asgn.student_id = u.id
		 JOIN quizzes q ON asgn.quiz_id = q.id
		 WHERE asgn.student_id = $1`,
		studentId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*model.Attempt
	for rows.Next() {
		a := &model.Attempt{User: &model.User{}, Quiz: &model.Quiz{}}
		if err = rows.Scan(
			&a.Id, &a.AssignmentId, &a.Score, &a.CorrectQuestionIds, &a.WrongQuestionIds,
			&a.User.Id, &a.User.Nickname,
			&a.Quiz.Id, &a.Quiz.Title,
		); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, nil
}
