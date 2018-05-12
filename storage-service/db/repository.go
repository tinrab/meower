package db

import "github.com/tinrab/meower/schema"

type Repository interface {
	Close() error
	InsertMeow(meow schema.Meow) error
	ListMeows(skip uint64, take uint64) ([]schema.Meow, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func InsertMeow(meow schema.Meow) error {
	return impl.InsertMeow(meow)
}

func ListMeows(skip uint64, take uint64) ([]schema.Meow, error) {
	return impl.ListMeows(skip, take)
}
