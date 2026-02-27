package dto

type AnswerOptionResponse struct {
	Id         int    `json:"id"`
	QuestionId int    `json:"question_id"`
	Text       string `json:"text"`
	Number     int    `json:"number"`
}

type QuestionResponse struct {
	Id                  int                    `json:"id"`
	QuizId              int                    `json:"quiz_id"`
	Text                string                 `json:"text"`
	CorrectAnswerNumber int                    `json:"correct_answer_number"`
	Options             []AnswerOptionResponse `json:"options"`
}

type QuizResponse struct {
	Id        int                `json:"id"`
	Title     string             `json:"title"`
	CreatorId int                `json:"creator_id"`
	Questions []QuestionResponse `json:"questions"`
}
