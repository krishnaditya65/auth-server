package http

import (
	"encoding/json"
	"net/http"
	"strings"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// isValidEmail performs a lightweight structural check: non-empty, contains exactly one @,
// with non-empty local and domain parts.
func isValidEmail(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return parts[0] != "" && parts[1] != "" && strings.Contains(parts[1], ".")
}
