package main

import (
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/nktch1/tank/internal/config"
	"github.com/nktch1/tank/internal/logger"
	"github.com/nktch1/tank/internal/middlware"
	"github.com/nktch1/tank/internal/server"
	"github.com/nktch1/tank/internal/service"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	var (
		lg   = logger.BuildLogger(conf)
		addr = fmt.Sprintf("0.0.0.0:%d", conf.Port)
		svc  = service.New(conf)
		srv  = server.New(svc)
	)

	srv = middlware.AddLogMiddleware(srv, lg)

	lg.Info("Server started at", zap.Int("port", conf.Port))

	if err = http.ListenAndServe(addr, srv); err != nil {
		log.Fatal(err)
	}
}
