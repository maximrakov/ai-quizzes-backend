package model

type AnswerOption struct {
	Id         int
	QuestionId int
	Text       string
	Number     int
}

type Question struct {
	Id                  int
	QuizId              int
	Text                string
	CorrectAnswerNumber int
	Options             []AnswerOption
}

type Quiz struct {
	Id        int
	Title     string
	CreatorId int
	Questions []Question
}

func NewQuiz(title string, creatorId int, questions []Question) *Quiz {
	return &Quiz{
		Title:     title,
		CreatorId: creatorId,
		Questions: questions,
	}
}
