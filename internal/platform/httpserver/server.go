package httpserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	port string
	mux  *chi.Mux
}

func New(
	port string,
) *Server {

	return &Server{
		port: port,
		mux:  chi.NewRouter(),
	}
}

func (s *Server) Router() *chi.Mux {
	return s.mux
}

func (s *Server) Start() error {

	addr := fmt.Sprintf(
		":%s",
		s.port,
	)

	return http.ListenAndServe(
		addr,
		s.mux,
	)
}
