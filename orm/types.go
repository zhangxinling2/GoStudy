package orm

import (
	"context"
)

type Queries[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMutil(ctx context.Context) ([]*T, error)
}

type Execute interface {
	Exec(ctx context.Context) error
}

type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
