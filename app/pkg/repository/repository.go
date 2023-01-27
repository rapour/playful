package repository

import (
	"context"
)

type ColumnRepository interface {
	Query(ctx context.Context, q string, values ...interface{}) error
}
