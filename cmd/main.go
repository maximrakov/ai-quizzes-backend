package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/maximrakov/ai-quizzes-backend/internal/app"
	"github.com/maximrakov/ai-quizzes-backend/internal/database/postgres"
	"github.com/maximrakov/ai-quizzes-backend/internal/handler"
	assignmentH "github.com/maximrakov/ai-quizzes-backend/internal/handler/assignment"
	attemptH "github.com/maximrakov/ai-quizzes-backend/internal/handler/attempt"
	quizH "github.com/maximrakov/ai-quizzes-backend/internal/handler/quiz"
	userH "github.com/maximrakov/ai-quizzes-backend/internal/handler/user"
	assignmentR "github.com/maximrakov/ai-quizzes-backend/internal/repository/assignment"
	attemptR "github.com/maximrakov/ai-quizzes-backend/internal/repository/attempt"
	quizR "github.com/maximrakov/ai-quizzes-backend/internal/repository/quiz"
	userR "github.com/maximrakov/ai-quizzes-backend/internal/repository/user"
	assignmentS "github.com/maximrakov/ai-quizzes-backend/internal/service/assignment"
	attemptS "github.com/maximrakov/ai-quizzes-backend/internal/service/attempt"
	quizS "github.com/maximrakov/ai-quizzes-backend/internal/service/quiz"
	userS "github.com/maximrakov/ai-quizzes-backend/internal/service/user"
	"github.com/maximrakov/ai-quizzes-backend/internal/static"
	pkgai "github.com/maximrakov/ai-quizzes-backend/pkg/ai"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env file: %v", err)
	}

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
	assignmentRepo := assignmentR.NewRepo(pg)
	attemptRepo := attemptR.NewRepo(pg)

	//ai client init
	aiClient, err := pkgai.NewOpenAIClient()
	if err != nil {
		appCtx.Logger.Error("failed to init openai client", "error", err)
		os.Exit(1)
	}

	//service init
	userService := userS.NewService(userRepo, appCtx.JwtSecret)
	quizService := quizS.NewService(quizRepo, userRepo, aiClient)
	assignmentService := assignmentS.NewService(assignmentRepo, userRepo)
	attemptService := attemptS.NewService(attemptRepo, assignmentRepo, quizRepo, userRepo)

	//handler init
	userHandler := userH.NewHandler(userService)
	quizHandler := quizH.NewHandler(quizService)
	assignmentHandler := assignmentH.NewHandler(assignmentService)
	attemptHandler := attemptH.NewHandler(attemptService)

	//init server
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, appCtx.JwtSecret, userHandler, quizHandler, assignmentHandler, attemptHandler)

	// serve embedded frontend
	webFS, err := fs.Sub(static.Web, "web")
	if err != nil {
		appCtx.Logger.Error("failed to create sub FS for frontend", "error", err)
		os.Exit(1)
	}
	mux.Handle("/", http.FileServerFS(webFS))

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler.LoggingMiddleware(appCtx.Logger, mux),
	}

	err = server.ListenAndServe()

	if err != nil {
		appCtx.Logger.Error("failed to start server", "error", err)
	}
}
