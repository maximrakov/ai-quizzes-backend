package handler

import "net/http"

type UserHandler interface {
	Create(http.ResponseWriter, *http.Request)
}

type QuizHandler interface {
	Create(http.ResponseWriter, *http.Request)
	Assign(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
	FindById(http.ResponseWriter, *http.Request)
	FindCreated(http.ResponseWriter, *http.Request)
	FindAssigned(http.ResponseWriter, *http.Request)
	FindByUserId(http.ResponseWriter, *http.Request)
}

func RegisterRoutes(mux *http.ServeMux, userHandler UserHandler, quizHandler QuizHandler) {
	mux.HandleFunc("POST /user", userHandler.Create)

	mux.HandleFunc("POST /quiz", quizHandler.Create)
	mux.HandleFunc("POST /quiz/{quizId}/assign", quizHandler.Assign)
	mux.HandleFunc("GET /quiz", quizHandler.FindAll)
	mux.HandleFunc("GET /quiz/{quizId}", quizHandler.FindById)
	mux.HandleFunc("GET /user/{userId}/quizzes", quizHandler.FindByUserId)
	mux.HandleFunc("GET /user/{userId}/quizzes/created", quizHandler.FindCreated)
	mux.HandleFunc("GET /user/{userId}/quizzes/assigned", quizHandler.FindAssigned)
}
