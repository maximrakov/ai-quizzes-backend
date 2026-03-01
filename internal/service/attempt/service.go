package attempt

import (
	"context"
	"errors"

	"github.com/maximrakov/ai-quizzes-backend/internal/model"
)

var ErrAssignmentNotForStudent = errors.New("this assignment does not belong to the given student")

type UserAnswer struct {
	QuestionId   int
	AnswerNumber int
}

type Repository interface {
	Create(ctx context.Context, attempt *model.Attempt) (*model.Attempt, error)
	FindByStudentId(ctx context.Context, studentId int) ([]*model.Attempt, error)
	FindByQuizId(ctx context.Context, quizId int) ([]*model.Attempt, error)
}

type AssignmentRepository interface {
	FindById(ctx context.Context, id int) (*model.Assignment, error)
}

type QuizRepository interface {
	FindById(ctx context.Context, id int) (*model.Quiz, error)
}

type UserRepository interface {
	FindById(ctx context.Context, id int) (*model.User, error)
}

type service struct {
	repo           Repository
	assignmentRepo AssignmentRepository
	quizRepo       QuizRepository
	userRepo       UserRepository
}

func NewService(repo Repository, assignmentRepo AssignmentRepository, quizRepo QuizRepository, userRepo UserRepository) *service {
	return &service{
		repo:           repo,
		assignmentRepo: assignmentRepo,
		quizRepo:       quizRepo,
		userRepo:       userRepo,
	}
}

func (s *service) GetByStudentId(ctx context.Context, studentId int) ([]*model.Attempt, error) {
	return s.repo.FindByStudentId(ctx, studentId)
}

func (s *service) GetByQuizId(ctx context.Context, quizId int) ([]*model.Attempt, error) {
	return s.repo.FindByQuizId(ctx, quizId)
}

func (s *service) Create(ctx context.Context, assignmentId int, answers []UserAnswer) (*model.Attempt, error) {
	assignment, err := s.assignmentRepo.FindById(ctx, assignmentId)
	if err != nil {
		return nil, err
	}

	quiz, err := s.quizRepo.FindById(ctx, assignment.QuizId)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindById(ctx, assignment.StudentId)
	if err != nil {
		return nil, err
	}

	answerMap := make(map[int]int, len(answers))
	for _, a := range answers {
		answerMap[a.QuestionId] = a.AnswerNumber
	}

	var correctIds, wrongIds []int32
	for _, q := range quiz.Questions {
		if answerMap[q.Id] == q.CorrectAnswerNumber {
			correctIds = append(correctIds, int32(q.Id))
		} else {
			wrongIds = append(wrongIds, int32(q.Id))
		}
	}

	if correctIds == nil {
		correctIds = []int32{}
	}
	if wrongIds == nil {
		wrongIds = []int32{}
	}

	var score float64
	if len(quiz.Questions) > 0 {
		score = float64(len(correctIds)) / float64(len(quiz.Questions)) * 100
	}

	attempt := model.NewAttempt(assignmentId, user, quiz, score, correctIds, wrongIds)
	return s.repo.Create(ctx, attempt)
}
