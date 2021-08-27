package service

import (
	"context"

	"github.com/nktch1/tank/internal/domain"
)

type Service interface {
	GatherHosts(ctx context.Context, req domain.Request) (*SearchResults, error)
	CheckResponsibility(ctx context.Context, hosts *SearchResults) (*domain.Response, error)
}