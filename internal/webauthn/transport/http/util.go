package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// requestWithBody returns a copy of r with its body replaced by the supplied bytes.
// The go-webauthn library reads the assertion from r.Body.
func requestWithBody(r *http.Request, body []byte) *http.Request {
	cp := r.Clone(r.Context())
	cp.Body = io.NopCloser(bytes.NewReader(body))
	cp.ContentLength = int64(len(body))
	return cp
}
