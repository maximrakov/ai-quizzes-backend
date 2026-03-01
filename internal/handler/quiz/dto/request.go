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
	Questions []QuestionRequest `json:"questions"`
}

type AssignQuizRequest struct {
	StudentId int `json:"student_id"`
	MentorId  int `json:"mentor_id"`
}

type GenerateQuestionsRequest struct {
	Title string `json:"title"`
	Topic string `json:"topic"`
	Count int    `json:"count"`
}
