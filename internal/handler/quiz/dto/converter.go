package dto

import "github.com/maximrakov/ai-quizzes-backend/internal/model"

func ToQuizResponse(quiz *model.Quiz) QuizResponse {
	return QuizResponse{
		Id:            quiz.Id,
		Question:      quiz.Question,
		Options:       quiz.Options,
		CorrectAnswer: quiz.CorrectAnswer,
		CreatorId:     quiz.CreatorId,
	}
}

func ToQuizResponses(quizzes []*model.Quiz) []QuizResponse {
	responses := make([]QuizResponse, 0, len(quizzes))
	for _, q := range quizzes {
		responses = append(responses, ToQuizResponse(q))
	}
	return responses
}
