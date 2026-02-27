package dto

import "github.com/maximrakov/ai-quizzes-backend/internal/model"

func ToAnswerOptionResponse(opt model.AnswerOption) AnswerOptionResponse {
	return AnswerOptionResponse{
		Id:         opt.Id,
		QuestionId: opt.QuestionId,
		Text:       opt.Text,
		Number:     opt.Number,
	}
}

func ToQuestionResponse(q model.Question) QuestionResponse {
	options := make([]AnswerOptionResponse, 0, len(q.Options))
	for _, opt := range q.Options {
		options = append(options, ToAnswerOptionResponse(opt))
	}
	return QuestionResponse{
		Id:                  q.Id,
		QuizId:              q.QuizId,
		Text:                q.Text,
		CorrectAnswerNumber: q.CorrectAnswerNumber,
		Options:             options,
	}
}

func ToQuizResponse(quiz *model.Quiz) QuizResponse {
	questions := make([]QuestionResponse, 0, len(quiz.Questions))
	for _, q := range quiz.Questions {
		questions = append(questions, ToQuestionResponse(q))
	}
	return QuizResponse{
		Id:        quiz.Id,
		Title:     quiz.Title,
		CreatorId: quiz.CreatorId,
		Questions: questions,
	}
}

func ToQuizResponses(quizzes []*model.Quiz) []QuizResponse {
	responses := make([]QuizResponse, 0, len(quizzes))
	for _, q := range quizzes {
		responses = append(responses, ToQuizResponse(q))
	}
	return responses
}

func ToQuestions(reqs []QuestionRequest) []model.Question {
	questions := make([]model.Question, 0, len(reqs))
	for _, qr := range reqs {
		options := make([]model.AnswerOption, 0, len(qr.Options))
		for _, or_ := range qr.Options {
			options = append(options, model.AnswerOption{
				Text:   or_.Text,
				Number: or_.Number,
			})
		}
		questions = append(questions, model.Question{
			Text:                qr.Text,
			CorrectAnswerNumber: qr.CorrectAnswerNumber,
			Options:             options,
		})
	}
	return questions
}
