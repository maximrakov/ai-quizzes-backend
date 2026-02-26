package app

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/maximrakov/ai-quizzes-backend/pkg/env"
)

type Context struct {
	Env         string
	Port        string
	PostgresUrl string
	Db          *sql.DB
	Logger      *slog.Logger
}

func NewContext() *Context {
	return &Context{
		Env:         env.GetEnv("env"),
		Port:        env.GetEnv("port"),
		PostgresUrl: os.Getenv("POSTGRES_URL"),
		Logger:      slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}
