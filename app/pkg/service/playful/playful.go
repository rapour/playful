package playful

import (
	"context"

	"playful/app/pkg/domain"
	"playful/app/pkg/repository"
	"playful/app/pkg/service"

	"github.com/gocql/gocql"
)

type playfulService struct {
	column repository.ColumnRepository
}

func (p *playfulService) SetLocation(ctx context.Context, loc domain.Location) error {

	return p.column.Query(ctx, `INSERT INTO location (id, altitude, longitude, time) VALUES (?, ?, ?, ?)`,
		gocql.TimeUUID(),
		loc.Altitude,
		loc.Longitude,
		loc.Timestamp,
	)

}

func NewPlayfulService(rep repository.ColumnRepository) (service.PlayfulService, error) {

	err := rep.Query(context.TODO(), "CREATE TABLE IF NOT EXISTS location (id UUID PRIMARY KEY, altitude int, longitude int, time timestamp)")
	if err != nil {
		return nil, err
	}

	return &playfulService{
		column: rep,
	}, nil

}
