package service

import (
	"context"
	"playful/app/pkg/domain"
)

type PlayfulService interface {
	SetLocation(ctx context.Context, loc domain.Location) error 
}