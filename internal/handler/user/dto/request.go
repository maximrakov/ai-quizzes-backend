package dto

import "github.com/maximrakov/ai-quizzes-backend/internal/model"

type RegisterRequest struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	Role     model.Role `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
