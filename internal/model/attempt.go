package model

type Attempt struct {
	Id                 int
	AssignmentId       int
	User               *User
	Quiz               *Quiz
	Score              float64
	CorrectQuestionIds []int32
	WrongQuestionIds   []int32
}

func NewAttempt(assignmentId int, user *User, quiz *Quiz, score float64, correctIds, wrongIds []int32) *Attempt {
	return &Attempt{
		AssignmentId:       assignmentId,
		User:               user,
		Quiz:               quiz,
		Score:              score,
		CorrectQuestionIds: correctIds,
		WrongQuestionIds:   wrongIds,
	}
}
