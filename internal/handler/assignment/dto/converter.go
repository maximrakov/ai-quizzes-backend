package dto

import "github.com/maximrakov/ai-quizzes-backend/internal/model"

func ToAssignmentResponse(a *model.Assignment) AssignmentResponse {
	return AssignmentResponse{
		Id:        a.Id,
		QuizId:    a.QuizId,
		StudentId: a.StudentId,
	}
}

func ToAssignmentResponses(assignments []*model.Assignment) []AssignmentResponse {
	responses := make([]AssignmentResponse, 0, len(assignments))
	for _, a := range assignments {
		responses = append(responses, ToAssignmentResponse(a))
	}
	return responses
}
