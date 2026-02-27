package main

import (
	"context"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/maximrakov/ai-quizzes-backend/internal/app"
	"github.com/maximrakov/ai-quizzes-backend/internal/database/postgres"
	"github.com/maximrakov/ai-quizzes-backend/internal/handler"
	quizH "github.com/maximrakov/ai-quizzes-backend/internal/handler/quiz"
	userH "github.com/maximrakov/ai-quizzes-backend/internal/handler/user"
	quizR "github.com/maximrakov/ai-quizzes-backend/internal/repository/quiz"
	userR "github.com/maximrakov/ai-quizzes-backend/internal/repository/user"
	quizS "github.com/maximrakov/ai-quizzes-backend/internal/service/quiz"
	userS "github.com/maximrakov/ai-quizzes-backend/internal/service/user"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	appCtx := app.NewContext()

	pg, err := postgres.New(*appCtx)

	if err != nil {
		appCtx.Logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	initializator := postgres.NewInitializer(pg)

	err = initializator.InitDB(ctx)

	if err != nil {
		appCtx.Logger.Error("failed to initialize PG")
	}

	appCtx.Logger.Info("PG inited")

	//repos init
	userRepo := userR.NewRepo(pg)
	quizRepo := quizR.NewRepo(pg)

	//service init
	userService := userS.NewService(userRepo)
	quizService := quizS.NewService(quizRepo, userRepo)

	//handler init
	userHandler := userH.NewHandler(userService)
	quizHandler := quizH.NewHandler(quizService)

	//init server
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, userHandler, quizHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()

	if err != nil {
		appCtx.Logger.Error("failed to start server", "error", err)
	}
}
