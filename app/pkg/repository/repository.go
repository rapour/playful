package repository

import (
	"context"

	"github.com/gocql/gocql"
)

type ColumnRepository interface {
	Get(ctx context.Context, q string, values ...interface{}) (int, gocql.Scanner)
	Set(ctx context.Context, q string, values ...interface{}) error
}
