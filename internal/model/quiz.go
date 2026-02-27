package model

type Quiz struct {
	Id            int
	Question      string
	Options       []string
	CorrectAnswer string
	CreatorId     int
}

func NewQuiz(question string, options []string, correctAnswer string, creatorId int) *Quiz {
	return &Quiz{
		Question:      question,
		Options:       options,
		CorrectAnswer: correctAnswer,
		CreatorId:     creatorId,
	}
}
