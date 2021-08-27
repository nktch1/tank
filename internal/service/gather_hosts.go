package service

import (
	"context"
	"sync"

	"github.com/nktch1/tank/internal/config"
	"github.com/nktch1/tank/internal/domain"
)

type Tank struct {
	Conf *config.Config

	sync.Mutex
}

func (t *Tank) GatherHosts(ctx context.Context, req domain.Request) (*SearchResults, error) {
	raw, err := Search(ctx, req.SearchString)
	if err != nil {
		return nil, err
	}

	return ParseSearchResults(ctx, raw), nil
}
