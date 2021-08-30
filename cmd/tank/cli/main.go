package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"github.com/nktch1/tank/internal/config"
	"github.com/nktch1/tank/internal/domain"
	"github.com/nktch1/tank/internal/logger"
	"github.com/nktch1/tank/internal/service"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	var (
		lg   = logger.BuildLogger(conf)
		svc  = service.New(conf)
		word string
	)

	if _, err = fmt.Scan(&word); err != nil {
		lg.Error("tried to parse hosts", zap.Error(err))
	}

	ctx := ctxzap.ToContext(context.Background(), lg)

	hosts, err := svc.GatherHosts(
		ctx,
		domain.Request{
			SearchString: word,
		},
	)

	if err != nil {
		lg.Error("tried to parse hosts", zap.Error(err))
	}

	resp, err := svc.CheckResponsibility(ctx, hosts)
	if err != nil {
		lg.Error("tried to check responsibility", zap.Error(err))
	}

	for k, v := range resp.HostToOptimalRPS {
		fmt.Println(k, v)
	}
}
