package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"github.com/nktch1/tank/internal/domain"
)

func (s server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logger := ctxzap.Extract(request.Context())

	if request.Method != http.MethodGet {
		logger.Error("method not allowed", zap.Int("status", http.StatusMethodNotAllowed))
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req domain.Request
	if err := schema.NewDecoder().Decode(&req, request.URL.Query()); err != nil {
		logger.Error("error", zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error("error", zap.Error(err))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debug("request",
		zap.String("search_string", req.SearchString),
	)

	hosts, err := s.Service.GatherHosts(request.Context(), req)
	if err != nil {
		logger.Error("error", zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := s.Service.CheckResponsibility(request.Context(), hosts)
	if err != nil {
		logger.Error("error", zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		logger.Error("error", zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = writer.Write(bytes); err != nil {
		logger.Error("error", zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
