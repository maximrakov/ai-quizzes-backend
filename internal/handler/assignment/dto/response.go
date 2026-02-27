package dto

type AssignmentResponse struct {
	Id        int `json:"id"`
	QuizId    int `json:"quiz_id"`
	StudentId int `json:"student_id"`
}
