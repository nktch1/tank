package service

import (
	"context"

	"github.com/nktch1/tank/internal/domain"
)

type Tank struct{}

func (t *Tank) CheckResponsibility(ctx context.Context, req domain.Request) (*domain.Response, error) {
	return nil, nil
}

