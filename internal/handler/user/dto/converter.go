package dto

import "github.com/maximrakov/ai-quizzes-backend/internal/model"

func ToUserResponse(user *model.User) UserResponse {
	return UserResponse{
		Id:       user.Id,
		Username: user.Nickname,
		Role:     string(user.Role),
	}
}

func ToUserResponses(users []*model.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = ToUserResponse(u)
	}
	return result
}
