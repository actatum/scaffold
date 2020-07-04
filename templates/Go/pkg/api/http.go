// package api

import (
	"net/http"
)

type Server interface {

}

type server struct {
	service service.Service
}

func NewServer(service service.Service) Server {
	return &server{
		service: service,
	}
}

func (s *server) HandleExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}