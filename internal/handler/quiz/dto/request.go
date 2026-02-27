package dto

type CreateQuizRequest struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correct_answer"`
	CreatorId     int      `json:"creator_id"`
}

type AssignQuizRequest struct {
	StudentId int `json:"student_id"`
	MentorId  int `json:"mentor_id"`
}
