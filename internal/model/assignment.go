package model

type Assignment struct {
	Id        int
	QuizId    int
	StudentId int
}

func NewAssignment(quizId, studentId int) *Assignment {
	return &Assignment{
		QuizId:    quizId,
		StudentId: studentId,
	}
}
