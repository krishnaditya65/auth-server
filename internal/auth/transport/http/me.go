package http

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) Me(
	w http.ResponseWriter,
	r *http.Request,
) {

	sessionID := r.Header.Get(
		"X-Session-ID",
	)

	if sessionID == "" {
		http.Error(
			w,
			"missing session",
			http.StatusUnauthorized,
		)
		return
	}

	result, err := h.meUseCase.Execute(
		r.Context(),
		sessionID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusUnauthorized,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		result,
	)
}
