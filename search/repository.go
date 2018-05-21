package search

import (
	"context"

	"github.com/tinrab/meower/schema"
)

type Repository interface {
	Close()
	InsertMeow(ctx context.Context, meow schema.Meow) error
	SearchMeows(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Meow, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertMeow(ctx context.Context, meow schema.Meow) error {
	return impl.InsertMeow(ctx, meow)
}

func SearchMeows(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Meow, error) {
	return impl.SearchMeows(ctx, query, skip, take)
}
