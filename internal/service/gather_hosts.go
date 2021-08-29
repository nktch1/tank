package service

import (
	"context"
	"net/http"
	"sync"

	"github.com/nktch1/tank/internal/config"
	"github.com/nktch1/tank/internal/domain"
)

type Tank struct {
	conf   *config.Config
	client http.Client

	sync.Mutex
}

func New(conf *config.Config) *Tank {
	return &Tank{
		conf: conf,
		client: http.Client{
			Timeout: conf.TimeoutPerHost,
		},
	}
}

func (t *Tank) GatherHosts(ctx context.Context, req domain.Request) (*SearchResults, error) {
	raw, err := Search(ctx, req.SearchString)
	if err != nil {
		return nil, err
	}

	return ParseSearchResults(ctx, raw), nil
}
