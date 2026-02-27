package quiz

import (
	"context"
	"errors"

	"github.com/maximrakov/ai-quizzes-backend/internal/model"
)

var (
	ErrNotMentor  = errors.New("only mentor can create quizzes")
	ErrNotStudent = errors.New("quiz can only be assigned to a student")
)

type Repository interface {
	Create(ctx context.Context, quiz *model.Quiz) (*model.Quiz, error)
	FindById(ctx context.Context, id int) (*model.Quiz, error)
	FindAll(ctx context.Context) ([]*model.Quiz, error)
	AssignToUser(ctx context.Context, quizId, userId int) error
	FindByUserId(ctx context.Context, userId int) ([]*model.Quiz, error)
	FindByCreatorId(ctx context.Context, creatorId int) ([]*model.Quiz, error)
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

func (s *service) Create(ctx context.Context, title string, creatorId int, questions []model.Question) (*model.Quiz, error) {
	creator, err := s.userRepo.FindById(ctx, creatorId)
	if err != nil {
		return nil, err
	}

	if creator.Role != model.Mentor {
		return nil, ErrNotMentor
	}

	quiz := model.NewQuiz(title, creatorId, questions)
	return s.repo.Create(ctx, quiz)
}

func (s *service) Assign(ctx context.Context, quizId, studentId, mentorId int) error {
	mentor, err := s.userRepo.FindById(ctx, mentorId)
	if err != nil {
		return err
	}
	if mentor.Role != model.Mentor {
		return ErrNotMentor
	}

	student, err := s.userRepo.FindById(ctx, studentId)
	if err != nil {
		return err
	}
	if student.Role != model.Student {
		return ErrNotStudent
	}

	return s.repo.AssignToUser(ctx, quizId, studentId)
}

func (s *service) FindAll(ctx context.Context) ([]*model.Quiz, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) FindById(ctx context.Context, id int) (*model.Quiz, error) {
	return s.repo.FindById(ctx, id)
}

func (s *service) FindByCreatorId(ctx context.Context, creatorId int) ([]*model.Quiz, error) {
	return s.repo.FindByCreatorId(ctx, creatorId)
}

func (s *service) FindByUserId(ctx context.Context, userId int) ([]*model.Quiz, error) {
	return s.repo.FindByUserId(ctx, userId)
}
