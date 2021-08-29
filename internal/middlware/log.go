package middlware

import (
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/nktch1/tank/internal/server"
)

type logMiddleware struct {
	next   server.Server
	logger *zap.Logger
}

func AddLogMiddleware(next server.Server, logger *zap.Logger) server.Server {
	return &logMiddleware{
		next:   next,
		logger: logger,
	}
}

func (m logMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	requestId := uuid.NewV4().String()

	request = request.WithContext(
		ctxzap.ToContext(request.Context(), m.logger.With(
			zap.String("request_id", requestId)),
		),
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
