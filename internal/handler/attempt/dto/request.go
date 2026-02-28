package dto

type UserAnswerRequest struct {
	QuestionId   int `json:"question_id"`
	AnswerNumber int `json:"answer_number"`
}

type CreateAttemptRequest struct {
	AssignmentId int                 `json:"assignment_id"`
	Answers      []UserAnswerRequest `json:"answers"`
}
