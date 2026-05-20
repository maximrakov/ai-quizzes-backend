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
	JwtSecret   string
	Db          *sql.DB
	Logger      *slog.Logger
}

func NewContext() *Context {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "changeme"
	}
	return &Context{
		Env:         env.GetEnv("env"),
		Port:        env.GetEnv("port"),
		PostgresUrl: os.Getenv("POSTGRES_URL"),
		JwtSecret:   secret,
		Logger:      slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}
