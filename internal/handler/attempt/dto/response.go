package dto

type AttemptUserResponse struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
}

type AttemptQuizResponse struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type AttemptResponse struct {
	Id                 int                 `json:"id"`
	User               AttemptUserResponse `json:"user"`
	Quiz               AttemptQuizResponse `json:"quiz"`
	Score              float64             `json:"score"`
	CorrectQuestionIds []int32             `json:"correct_question_ids"`
	WrongQuestionIds   []int32             `json:"wrong_question_ids"`
}
