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

func New(port string) *Server {
	return &Server{
		port: port,
		mux:  chi.NewRouter(),
	}
}

func (s *Server) Router() *chi.Mux {
	return s.mux
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	return http.ListenAndServe(addr, s.mux)
}

// CORS is a chi-compatible middleware that adds permissive CORS headers and handles preflight requests.
// Restrict Access-Control-Allow-Origin in production once a frontend origin is known.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Session-ID")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
