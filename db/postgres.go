package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/tinrab/meower/schema"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(connection string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) InsertMeow(meow schema.Meow) error {
	_, err := r.db.Exec("INSERT INTO meows(id, body) VALUES($1, $2)", meow.ID, meow.Body)
	return err
}
