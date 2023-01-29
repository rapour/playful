package playful

import (
	"context"
	"fmt"
	"time"

	"playful/app/pkg/domain"
	"playful/app/pkg/repository"
	"playful/app/pkg/service"
)

type playfulService struct {
	column repository.ColumnRepository
}

func (p *playfulService) SetLocation(ctx context.Context, loc domain.Location) error {

	err := p.column.Set(ctx, `INSERT INTO locations (id, time, altitude, longitude) VALUES (?, ?, ?, ?);`,
		loc.Ident,
		time.Unix(int64(loc.Timestamp), 0),
		loc.Altitude,
		loc.Longitude,
	)

	if err != nil {
		return fmt.Errorf("[SetLocation]: %v", err)
	}

	return nil

}

func (p *playfulService) GetLoaction(ctx context.Context) ([]domain.Location, error) {

	num, scanner := p.column.Get(ctx, "SELECT id, altitude, longitude, time FROM locations LIMIT 1")

	var time time.Time
	var result domain.Location
	var resultSet []domain.Location
	for scanner.Next() {
		err := scanner.Scan(&result.Ident, &result.Altitude, &result.Longitude, &time)
		result.Timestamp = int32(time.Unix())

		if err != nil {
			return nil, fmt.Errorf("[GetLocation][rows: %d]: %v", num, err)
		}
		resultSet = append(resultSet, result)
	}
	if len(resultSet) == 0 {
		return nil, fmt.Errorf("[GetLocation][rows: %d]: %v", num, scanner.Err())
	}

	return resultSet, nil
}

func NewPlayfulService(rep repository.ColumnRepository) (service.PlayfulService, error) {

	err := rep.Set(context.TODO(), "CREATE TABLE IF NOT EXISTS locations (id int, time timestamp, altitude int, longitude int, PRIMARY KEY (id, time)) WITH CLUSTERING ORDER BY (time DESC)")
	if err != nil {
		return nil, err
	}

	return &playfulService{
		column: rep,
	}, nil

}
