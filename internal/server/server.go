package server

import (
	"encoding/json"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"

	"github.com/nktch1/tank/internal/domain"
	"github.com/nktch1/tank/internal/service"
)

type Server interface {
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}

type server struct {
	Service service.Service
}

func New(svc service.Service) Server {
	return server{svc}
}

func (s server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	lg := ctxzap.Extract(request.Context())

	if request.Method != http.MethodGet {
		lg.Error("method not allowed", zap.Int("status", http.StatusMethodNotAllowed))
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var req domain.Request
	if err = json.Unmarshal(body, &req); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := s.Service.CheckResponsibility(request.Context(), req)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = writer.Write(bytes); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
