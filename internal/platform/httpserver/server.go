package httpserver

import (
	"fmt"
	"net/http"
)

type Server struct {
	port string
	mux  *http.ServeMux
}

func New(port string) *Server {
	return &Server{
		port: port,
		mux:  http.NewServeMux(),
	}
}

func (s *Server) Handle(
	method string,
	path string,
	handler http.Handler,
) {

	s.mux.Handle(
		path,
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if r.Method != method {
					http.Error(
						w,
						"method not allowed",
						http.StatusMethodNotAllowed,
					)
					return
				}

				handler.ServeHTTP(
					w,
					r,
				)
			},
		),
	)
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
