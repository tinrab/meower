package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/tinrab/meower/schema"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) InsertMeow(ctx context.Context, meow schema.Meow) error {
	_, err := r.db.Exec("INSERT INTO meows(id, body, created_at) VALUES($1, $2, $3)", meow.ID, meow.Body, meow.CreatedAt)
	return err
}

func (r *PostgresRepository) ListMeows(ctx context.Context, skip uint64, take uint64) ([]schema.Meow, error) {
	rows, err := r.db.Query("SELECT * FROM meows ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse all rows into an array of Meows
	meows := []schema.Meow{}
	for rows.Next() {
		meow := schema.Meow{}
		if err = rows.Scan(&meow.ID, &meow.Body, &meow.CreatedAt); err == nil {
			meows = append(meows, meow)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return meows, nil
}
