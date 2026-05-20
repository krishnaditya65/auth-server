package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	port string
}

func New(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			log.Printf("write response error: %v", err)
		}
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server.ListenAndServe()
}
