package dto

import "github.com/maximrakov/ai-quizzes-backend/internal/model"

func ToAttemptResponse(a *model.Attempt) AttemptResponse {
	return AttemptResponse{
		Id: a.Id,
		User: AttemptUserResponse{
			Id:       a.User.Id,
			Nickname: a.User.Nickname,
		},
		Quiz: AttemptQuizResponse{
			Id:    a.Quiz.Id,
			Title: a.Quiz.Title,
		},
		Score:              a.Score,
		CorrectQuestionIds: a.CorrectQuestionIds,
		WrongQuestionIds:   a.WrongQuestionIds,
	}
}
