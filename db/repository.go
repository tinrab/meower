package db

import "github.com/tinrab/meower/schema"

type Repository interface {
	Close() error
	InsertMeow(meow schema.Meow) error
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func InsertMeow(meow schema.Meow) error {
	return impl.InsertMeow(meow)
}
