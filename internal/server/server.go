package server

import (
	"net/http"

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