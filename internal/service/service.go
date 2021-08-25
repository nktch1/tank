package service

import (
	"context"
	"github.com/nktch1/tank/internal/domain"
)

type Service interface {
	CheckResponsibility(ctx context.Context, req domain.Request) (*domain.Response, error)
}
