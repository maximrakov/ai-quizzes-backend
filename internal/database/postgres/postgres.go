package postgres

import (
	"context"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximrakov/ai-quizzes-backend/internal/app"
)

type Postgres struct {
	Pool          *pgxpool.Pool
	DurationQuery time.Duration
	nameDb        string
	hostDb        string
}

func New(ctx app.Context) (*Postgres, error) {
	postgres := &Postgres{}

	postgres.DurationQuery = 300 * time.Second

	config, err := pgxpool.ParseConfig(ctx.PostgresUrl)

	if err != nil {
		ctx.Logger.Error("failed to parse postgres url", "url", ctx.PostgresUrl, "error", err)
		return nil, err
	}

	config.MinConns = 1
	config.MaxConns = 3

	postgres.Pool, err = pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		ctx.Logger.Error("failed to connect to postgres", "error", err)
		return nil, err
	}

	err = postgres.Ping(context.Background())

	if err != nil {
		ctx.Logger.Error("failed to ping postgres", "error", err)
		return nil, err
	}

	reDbName := regexp.MustCompile(`dbname=([^ ]+)`)
	reDbHost := regexp.MustCompile(`host=([^ ]+)`)

	if matches := reDbName.FindStringSubmatch(ctx.PostgresUrl); len(matches) > 1 {
		postgres.nameDb = matches[1]
	}
	if matches := reDbHost.FindStringSubmatch(ctx.PostgresUrl); len(matches) > 1 {
		postgres.hostDb = matches[1]
	}

	ctx.Logger.Info("connected to postgres", "DbName", postgres.nameDb, "DbHost", postgres.hostDb)

	return postgres, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.Pool.Ping(ctx)
}
