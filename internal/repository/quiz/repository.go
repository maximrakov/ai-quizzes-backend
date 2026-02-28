package quiz

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

func (r *repository) Create(ctx context.Context, quiz *model.Quiz) (*model.Quiz, error) {
	var quizId int
	err := r.postgres.Pool.QueryRow(ctx,
		"INSERT INTO quizzes (title, creator_id) VALUES ($1, $2) RETURNING id",
		quiz.Title, quiz.CreatorId,
	).Scan(&quizId)
	if err != nil {
		return nil, err
	}
	quiz.Id = quizId

	for i := range quiz.Questions {
		q := &quiz.Questions[i]
		q.QuizId = quizId

		var questionId int
		err = r.postgres.Pool.QueryRow(ctx,
			"INSERT INTO questions (quiz_id, text, correct_answer_number) VALUES ($1, $2, $3) RETURNING id",
			quizId, q.Text, q.CorrectAnswerNumber,
		).Scan(&questionId)
		if err != nil {
			return nil, err
		}
		q.Id = questionId

		for j := range q.Options {
			opt := &q.Options[j]
			opt.QuestionId = questionId

			var optId int
			err = r.postgres.Pool.QueryRow(ctx,
				"INSERT INTO answer_options (question_id, text, number) VALUES ($1, $2, $3) RETURNING id",
				questionId, opt.Text, opt.Number,
			).Scan(&optId)
			if err != nil {
				return nil, err
			}
			opt.Id = optId
		}
	}

	return quiz, nil
}

func (r *repository) FindById(ctx context.Context, id int) (*model.Quiz, error) {
	quiz := &model.Quiz{}
	err := r.postgres.Pool.QueryRow(ctx,
		"SELECT id, title, creator_id FROM quizzes WHERE id = $1", id,
	).Scan(&quiz.Id, &quiz.Title, &quiz.CreatorId)
	if err != nil {
		return nil, err
	}

	if err = r.loadQuestions(ctx, quiz); err != nil {
		return nil, err
	}

	return quiz, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*model.Quiz, error) {
	rows, err := r.postgres.Pool.Query(ctx,
		"SELECT id, title, creator_id FROM quizzes",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []*model.Quiz
	for rows.Next() {
		quiz := &model.Quiz{}
		if err = rows.Scan(&quiz.Id, &quiz.Title, &quiz.CreatorId); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	rows.Close()

	for _, quiz := range quizzes {
		if err = r.loadQuestions(ctx, quiz); err != nil {
			return nil, err
		}
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
		"SELECT id, title, creator_id FROM quizzes WHERE creator_id = $1", creatorId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []*model.Quiz
	for rows.Next() {
		quiz := &model.Quiz{}
		if err = rows.Scan(&quiz.Id, &quiz.Title, &quiz.CreatorId); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	rows.Close()

	for _, quiz := range quizzes {
		if err = r.loadQuestions(ctx, quiz); err != nil {
			return nil, err
		}
	}

	return quizzes, nil
}

func (r *repository) FindByUserId(ctx context.Context, userId int) ([]*model.Quiz, error) {
	rows, err := r.postgres.Pool.Query(ctx,
		`SELECT q.id, q.title, q.creator_id
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
		if err = rows.Scan(&quiz.Id, &quiz.Title, &quiz.CreatorId); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	rows.Close()

	for _, quiz := range quizzes {
		if err = r.loadQuestions(ctx, quiz); err != nil {
			return nil, err
		}
	}

	return quizzes, nil
}

func (r *repository) loadQuestions(ctx context.Context, quiz *model.Quiz) error {
	rows, err := r.postgres.Pool.Query(ctx,
		"SELECT id, quiz_id, text, correct_answer_number FROM questions WHERE quiz_id = $1 ORDER BY id",
		quiz.Id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var questions []model.Question
	for rows.Next() {
		var q model.Question
		if err = rows.Scan(&q.Id, &q.QuizId, &q.Text, &q.CorrectAnswerNumber); err != nil {
			return err
		}
		questions = append(questions, q)
	}
	rows.Close()

	for i := range questions {
		if err = r.loadOptions(ctx, &questions[i]); err != nil {
			return err
		}
	}

	quiz.Questions = questions
	return nil
}

func (r *repository) loadOptions(ctx context.Context, question *model.Question) error {
	rows, err := r.postgres.Pool.Query(ctx,
		"SELECT id, question_id, text, number FROM answer_options WHERE question_id = $1 ORDER BY number",
		question.Id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var options []model.AnswerOption
	for rows.Next() {
		var opt model.AnswerOption
		if err = rows.Scan(&opt.Id, &opt.QuestionId, &opt.Text, &opt.Number); err != nil {
			return err
		}
		options = append(options, opt)
	}

	question.Options = options
	return nil
}
