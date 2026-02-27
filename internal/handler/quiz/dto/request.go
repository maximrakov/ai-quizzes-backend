package dto

type AnswerOptionRequest struct {
	Text   string `json:"text"`
	Number int    `json:"number"`
}

type QuestionRequest struct {
	Text                string                `json:"text"`
	CorrectAnswerNumber int                   `json:"correct_answer_number"`
	Options             []AnswerOptionRequest `json:"options"`
}

type CreateQuizRequest struct {
	Title     string            `json:"title"`
	CreatorId int               `json:"creator_id"`
	Questions []QuestionRequest `json:"questions"`
}

type AssignQuizRequest struct {
	StudentId int `json:"student_id"`
	MentorId  int `json:"mentor_id"`
}
