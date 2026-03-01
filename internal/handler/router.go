package handler

import "net/http"

type UserHandler interface {
	Create(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
}

type QuizHandler interface {
	Create(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
	FindById(http.ResponseWriter, *http.Request)
	FindCreated(http.ResponseWriter, *http.Request)
	FindAssigned(http.ResponseWriter, *http.Request)
	FindByUserId(http.ResponseWriter, *http.Request)
	Generate(http.ResponseWriter, *http.Request)
}

type AssignmentHandler interface {
	Create(http.ResponseWriter, *http.Request)
	FindById(http.ResponseWriter, *http.Request)
	FindByStudentId(http.ResponseWriter, *http.Request)
	FindByQuizId(http.ResponseWriter, *http.Request)
}

type AttemptHandler interface {
	Create(http.ResponseWriter, *http.Request)
}

func auth(secret string, h http.HandlerFunc) http.HandlerFunc {
	return AuthMiddleware(secret, h).ServeHTTP
}

func RegisterRoutes(mux *http.ServeMux, jwtSecret string, userHandler UserHandler, quizHandler QuizHandler, assignmentHandler AssignmentHandler, attemptHandler AttemptHandler) {
	// публичные маршруты
	mux.HandleFunc("POST /user", userHandler.Create)
	mux.HandleFunc("POST /auth/login", userHandler.Login)

	mux.HandleFunc("GET /user", auth(jwtSecret, userHandler.FindAll))

	// защищённые маршруты
	mux.HandleFunc("POST /quiz", auth(jwtSecret, quizHandler.Create))
	mux.HandleFunc("POST /quiz/generate", auth(jwtSecret, quizHandler.Generate))
	mux.HandleFunc("GET /quiz", auth(jwtSecret, quizHandler.FindAll))
	mux.HandleFunc("GET /quiz/{quizId}", auth(jwtSecret, quizHandler.FindById))
	mux.HandleFunc("GET /user/{userId}/quizzes", auth(jwtSecret, quizHandler.FindByUserId))
	mux.HandleFunc("GET /user/{userId}/quizzes/created", auth(jwtSecret, quizHandler.FindCreated))
	mux.HandleFunc("GET /user/{userId}/quizzes/assigned", auth(jwtSecret, quizHandler.FindAssigned))

	mux.HandleFunc("POST /assignment", auth(jwtSecret, assignmentHandler.Create))
	mux.HandleFunc("GET /assignment/{id}", auth(jwtSecret, assignmentHandler.FindById))
	mux.HandleFunc("GET /user/{userId}/assignments", auth(jwtSecret, assignmentHandler.FindByStudentId))
	mux.HandleFunc("GET /quiz/{quizId}/assignments", auth(jwtSecret, assignmentHandler.FindByQuizId))

	mux.HandleFunc("POST /attempt", auth(jwtSecret, attemptHandler.Create))
}
