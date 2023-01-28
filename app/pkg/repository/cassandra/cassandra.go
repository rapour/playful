package cassandra

import (
	"context"

	"playful/app/pkg/repository"
	cassandra_tools "playful/app/tools/cassandra"

	"github.com/gocql/gocql"
)

type cassandraRepository struct {
	client *cassandra_tools.Client
}

func (r *cassandraRepository) Set(ctx context.Context, q string, values ...interface{}) error {

	query := r.client.Session.Query(q, values...).WithContext(ctx)
	return query.Exec()
}

func (r *cassandraRepository) Get(ctx context.Context, q string, values ...interface{}) (int, gocql.Scanner) {

	iter := r.client.Session.Query(q, values...).WithContext(ctx).Iter()
	return iter.NumRows(), iter.Scanner()
}

func NewColumnRepository(cli *cassandra_tools.Client) repository.ColumnRepository {

	return &cassandraRepository{
		client: cli,
	}
}
