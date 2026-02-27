package dto

type CreateAssignmentRequest struct {
	QuizId    int `json:"quiz_id"`
	StudentId int `json:"student_id"`
	MentorId  int `json:"mentor_id"`
}
