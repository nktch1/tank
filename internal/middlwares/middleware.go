package middlwares

import (
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"github.com/nktch1/tank/internal/server"
)

type middleware struct {
	next   server.Server
	logger *zap.Logger
}

func AddMiddleware(next server.Server, logger *zap.Logger) server.Server {
	return &middleware{
		next:   next,
		logger: logger,
	}
}

func (m middleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	request = request.WithContext(
		ctxzap.ToContext(request.Context(), m.logger),
	)

	start := time.Now()

	m.next.ServeHTTP(writer, request)

	ctxzap.Extract(
		request.Context()).Info(
			"finished",
			zap.Int64("time_ms",
			time.Since(start).Milliseconds(),
		),
	)
}
