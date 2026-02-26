package main

import (
	"context"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/maximrakov/ai-quizzes-backend/internal/app"
	"github.com/maximrakov/ai-quizzes-backend/internal/database/postgres"
	"github.com/maximrakov/ai-quizzes-backend/internal/handler"
	userH "github.com/maximrakov/ai-quizzes-backend/internal/handler/user"
	userR "github.com/maximrakov/ai-quizzes-backend/internal/repository/user"
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

	//service init
	userService := userS.NewService(userRepo)

	//handler init
	userHandler := userH.NewHandler(userService)

	//init server
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, userHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()

	if err != nil {
		appCtx.Logger.Error("failed to start server", "error", err)
	}
}
