package postgres

import (
	"context"
	"fmt"
)

type Initializator struct {
	pg *Postgres
}

func NewInitializer(pg *Postgres) *Initializator {
	return &Initializator{pg: pg}
}

func (i *Initializator) InitDB(ctx context.Context) error {
	initQuery := []string{
		`CREATE TABLE IF NOT EXISTS users (
    		id serial PRIMARY KEY,
    		nickname VARCHAR(255) NOT NULL,
    		password VARCHAR(255) NOT NULL,
    		role VARCHAR(127) NOT NULL
		);
		`,
		`CREATE TABLE IF NOT EXISTS quizzes (
			id SERIAL PRIMARY KEY,
			question TEXT NOT NULL,
			options JSONB NOT NULL,
			correct_answer VARCHAR(255) NOT NULL,
			creator_id INTEGER NOT NULL REFERENCES users(id)
		);
		`,
		`CREATE TABLE IF NOT EXISTS user_quizzes (
			user_id INTEGER NOT NULL REFERENCES users(id),
			quiz_id INTEGER NOT NULL REFERENCES quizzes(id),
			PRIMARY KEY (user_id, quiz_id)
		);
		`,
	}

	for _, query := range initQuery {
		ctx, cancel := context.WithTimeout(ctx, i.pg.DurationQuery)
		defer cancel()
		res, err := i.pg.Pool.Exec(ctx, query)
		if err != nil {
			return err
		}
		fmt.Println(res.String())
	}

	return nil
}
