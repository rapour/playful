package cassandra

import (
	"context"

	"playful/app/pkg/repository"
	cassandra_tools "playful/app/tools/cassandra"
)

type cassandraRepository struct {
	client *cassandra_tools.Client
}

func (r *cassandraRepository) Query(ctx context.Context, q string, values ...interface{}) error {

	query := r.client.Session.Query(q, values...).WithContext(ctx)
	return query.Exec()
}

func NewColumnRepository(cli *cassandra_tools.Client) repository.ColumnRepository {

	return &cassandraRepository{
		client: cli,
	}
}
