package handler

import "net/http"

type UserHandler interface {
	Create(http.ResponseWriter, *http.Request)
}

type QuizHandler interface {
	Create(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
	FindById(http.ResponseWriter, *http.Request)
	FindCreated(http.ResponseWriter, *http.Request)
	FindAssigned(http.ResponseWriter, *http.Request)
	FindByUserId(http.ResponseWriter, *http.Request)
}

type AssignmentHandler interface {
	Create(http.ResponseWriter, *http.Request)
	FindById(http.ResponseWriter, *http.Request)
	FindByStudentId(http.ResponseWriter, *http.Request)
}

type AttemptHandler interface {
	Create(http.ResponseWriter, *http.Request)
}

func RegisterRoutes(mux *http.ServeMux, userHandler UserHandler, quizHandler QuizHandler, assignmentHandler AssignmentHandler, attemptHandler AttemptHandler) {
	mux.HandleFunc("POST /user", userHandler.Create)

	mux.HandleFunc("POST /quiz", quizHandler.Create)
	mux.HandleFunc("GET /quiz", quizHandler.FindAll)
	mux.HandleFunc("GET /quiz/{quizId}", quizHandler.FindById)
	mux.HandleFunc("GET /user/{userId}/quizzes", quizHandler.FindByUserId)
	mux.HandleFunc("GET /user/{userId}/quizzes/created", quizHandler.FindCreated)
	mux.HandleFunc("GET /user/{userId}/quizzes/assigned", quizHandler.FindAssigned)

	mux.HandleFunc("POST /assignment", assignmentHandler.Create)
	mux.HandleFunc("GET /assignment/{id}", assignmentHandler.FindById)
	mux.HandleFunc("GET /user/{userId}/assignments", assignmentHandler.FindByStudentId)

	mux.HandleFunc("POST /attempt", attemptHandler.Create)
}
