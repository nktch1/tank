package main

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"

	"github.com/nktch1/tank/internal/config"
	"github.com/nktch1/tank/internal/logging"
	"github.com/nktch1/tank/internal/middlwares"
	"github.com/nktch1/tank/internal/server"
	"github.com/nktch1/tank/internal/service"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	var (
		logger = logging.BuildLogger(conf)
		addr   = fmt.Sprintf("0.0.0.0:%d", conf.Port)
		svc    = &service.Tank{}
		srv    = server.New(svc)
	)

	srv = middlwares.AddMiddleware(srv, logger)

	logger.Info("Server started at", zap.Int("port", conf.Port))

	if err = http.ListenAndServe(addr, srv); err != nil {
		log.Fatal(err)
	}
}
