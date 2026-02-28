package user

import (
	"context"
	"fmt"

	"github.com/maximrakov/ai-quizzes-backend/internal/database/postgres"
	"github.com/maximrakov/ai-quizzes-backend/internal/model"
)

type repository struct {
	postgres *postgres.Postgres
}

func NewRepo(postgres *postgres.Postgres) *repository {
	return &repository{postgres: postgres}
}

func (r *repository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	var userId int
	err := r.postgres.Pool.QueryRow(ctx, "INSERT INTO USERS (nickname, password, role) VALUES ($1, $2, $3) RETURNING id", user.Nickname, user.Password, user.Role).
		Scan(&userId)

	if err != nil {
		return nil, err
	}

	user.Id = userId
	fmt.Println(user)
	return user, nil
}

func (r *repository) FindById(ctx context.Context, id int) (*model.User, error) {
	user := &model.User{}
	err := r.postgres.Pool.QueryRow(ctx, "SELECT id, nickname, password, role FROM users WHERE id = $1", id).
		Scan(&user.Id, &user.Nickname, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) FindByNickname(ctx context.Context, nickname string) (*model.User, error) {
	user := &model.User{}
	err := r.postgres.Pool.QueryRow(ctx, "SELECT id, nickname, password, role FROM users WHERE nickname = $1", nickname).
		Scan(&user.Id, &user.Nickname, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}

	return user, nil
}
