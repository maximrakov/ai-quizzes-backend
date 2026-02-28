package dto

import "github.com/maximrakov/ai-quizzes-backend/internal/model"

func ToUserResponse(user *model.User) UserResponse {

	return UserResponse{
		Id:       user.Id,
		Username: user.Nickname,
		Role:     string(user.Role),
	}
}
