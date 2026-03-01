package assignment

import (
	"context"
	"errors"

	"github.com/maximrakov/ai-quizzes-backend/internal/model"
)

var (
	ErrNotMentor  = errors.New("only mentor can assign quizzes")
	ErrNotStudent = errors.New("quiz can only be assigned to a student")
)

type Repository interface {
	Create(ctx context.Context, assignment *model.Assignment) (*model.Assignment, error)
	FindById(ctx context.Context, id int) (*model.Assignment, error)
	FindByStudentId(ctx context.Context, studentId int) ([]*model.Assignment, error)
	FindByQuizId(ctx context.Context, quizId int) ([]*model.Assignment, error)
}

type UserRepository interface {
	FindById(ctx context.Context, id int) (*model.User, error)
}

type service struct {
	repo     Repository
	userRepo UserRepository
}

func NewService(repo Repository, userRepo UserRepository) *service {
	return &service{repo: repo, userRepo: userRepo}
}

func (s *service) Create(ctx context.Context, quizId, studentId, mentorId int) (*model.Assignment, error) {
	mentor, err := s.userRepo.FindById(ctx, mentorId)
	if err != nil {
		return nil, err
	}
	if mentor.Role != model.Mentor {
		return nil, ErrNotMentor
	}

	student, err := s.userRepo.FindById(ctx, studentId)
	if err != nil {
		return nil, err
	}
	if student.Role != model.Student {
		return nil, ErrNotStudent
	}

	return s.repo.Create(ctx, model.NewAssignment(quizId, studentId))
}

func (s *service) FindById(ctx context.Context, id int) (*model.Assignment, error) {
	return s.repo.FindById(ctx, id)
}

func (s *service) FindByStudentId(ctx context.Context, studentId int) ([]*model.Assignment, error) {
	return s.repo.FindByStudentId(ctx, studentId)
}

func (s *service) FindByQuizId(ctx context.Context, quizId int) ([]*model.Assignment, error) {
	return s.repo.FindByQuizId(ctx, quizId)
}
