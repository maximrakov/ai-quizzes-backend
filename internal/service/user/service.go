package user

import (
	"context"
	"errors"

	"github.com/maximrakov/ai-quizzes-backend/internal/model"
	pkgjwt "github.com/maximrakov/ai-quizzes-backend/pkg/jwt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Repository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindByNickname(ctx context.Context, nickname string) (*model.User, error)
	FindAll(ctx context.Context, role string) ([]*model.User, error)
}

type service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) *service {
	return &service{repo: repo, jwtSecret: jwtSecret}
}

func (s *service) Create(ctx context.Context, nickname, password string, role model.Role) (*model.User, error) {
	user := model.NewUser(nickname, password, role)
	user, err := s.repo.Create(ctx, user)
	return user, err
}

func (s *service) FindAll(ctx context.Context, role string) ([]*model.User, error) {
	return s.repo.FindAll(ctx, role)
}

func (s *service) Login(ctx context.Context, nickname, password string) (string, error) {
	user, err := s.repo.FindByNickname(ctx, nickname)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if user.Password != password {
		return "", ErrInvalidCredentials
	}

	token, err := pkgjwt.Generate(user.Id, string(user.Role), s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
