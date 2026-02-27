package quiz

import (
	"context"
	"encoding/json"

	"github.com/maximrakov/ai-quizzes-backend/internal/database/postgres"
	"github.com/maximrakov/ai-quizzes-backend/internal/model"
)

type repository struct {
	postgres *postgres.Postgres
}

func NewRepo(postgres *postgres.Postgres) *repository {
	return &repository{postgres: postgres}
}

func (r *repository) Create(ctx context.Context, quiz *model.Quiz) (*model.Quiz, error) {
	optionsJSON, err := json.Marshal(quiz.Options)
	if err != nil {
		return nil, err
	}

	var quizId int
	err = r.postgres.Pool.QueryRow(ctx,
		"INSERT INTO quizzes (question, options, correct_answer, creator_id) VALUES ($1, $2, $3, $4) RETURNING id",
		quiz.Question, optionsJSON, quiz.CorrectAnswer, quiz.CreatorId,
	).Scan(&quizId)

	if err != nil {
		return nil, err
	}

	quiz.Id = quizId
	return quiz, nil
}

func (r *repository) FindById(ctx context.Context, id int) (*model.Quiz, error) {
	quiz := &model.Quiz{}
	var optionsRaw []byte

	err := r.postgres.Pool.QueryRow(ctx,
		"SELECT id, question, options, correct_answer, creator_id FROM quizzes WHERE id = $1", id,
	).Scan(&quiz.Id, &quiz.Question, &optionsRaw, &quiz.CorrectAnswer, &quiz.CreatorId)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(optionsRaw, &quiz.Options); err != nil {
		return nil, err
	}

	return quiz, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*model.Quiz, error) {
	rows, err := r.postgres.Pool.Query(ctx,
		"SELECT id, question, options, correct_answer, creator_id FROM quizzes",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []*model.Quiz
	for rows.Next() {
		quiz := &model.Quiz{}
		var optionsRaw []byte

		if err = rows.Scan(&quiz.Id, &quiz.Question, &optionsRaw, &quiz.CorrectAnswer, &quiz.CreatorId); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(optionsRaw, &quiz.Options); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

func (r *repository) AssignToUser(ctx context.Context, quizId, userId int) error {
	_, err := r.postgres.Pool.Exec(ctx,
		"INSERT INTO user_quizzes (user_id, quiz_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		userId, quizId,
	)
	return err
}

func (r *repository) FindByCreatorId(ctx context.Context, creatorId int) ([]*model.Quiz, error) {
	rows, err := r.postgres.Pool.Query(ctx,
		"SELECT id, question, options, correct_answer, creator_id FROM quizzes WHERE creator_id = $1",
		creatorId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []*model.Quiz
	for rows.Next() {
		quiz := &model.Quiz{}
		var optionsRaw []byte

		if err = rows.Scan(&quiz.Id, &quiz.Question, &optionsRaw, &quiz.CorrectAnswer, &quiz.CreatorId); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(optionsRaw, &quiz.Options); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

func (r *repository) FindByUserId(ctx context.Context, userId int) ([]*model.Quiz, error) {
	rows, err := r.postgres.Pool.Query(ctx,
		`SELECT q.id, q.question, q.options, q.correct_answer, q.creator_id
		 FROM quizzes q
		 JOIN user_quizzes uq ON uq.quiz_id = q.id
		 WHERE uq.user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []*model.Quiz
	for rows.Next() {
		quiz := &model.Quiz{}
		var optionsRaw []byte

		if err = rows.Scan(&quiz.Id, &quiz.Question, &optionsRaw, &quiz.CorrectAnswer, &quiz.CreatorId); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(optionsRaw, &quiz.Options); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}
