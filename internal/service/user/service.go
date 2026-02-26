package user

import (
	"context"

	"github.com/maximrakov/ai-quizzes-backend/internal/model"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
}
type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, nickname, password string, role model.Role) (*model.User, error) {
	user := model.NewUser(nickname, password, role)
	user, err := s.repo.Create(ctx, user)

	return user, err
}
