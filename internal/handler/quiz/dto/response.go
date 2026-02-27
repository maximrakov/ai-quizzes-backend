package dto

type QuizResponse struct {
	Id            int      `json:"id"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correct_answer"`
	CreatorId     int      `json:"creator_id"`
}
