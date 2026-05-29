package http

import (
	"encoding/json"
	"net/http"

	authctx "github.com/krishnaditya65/auth-server/internal/shared/context"
)

func (h *Handler) Me(
	w http.ResponseWriter,
	r *http.Request,
) {
	p := authctx.MustPrincipal(
		r.Context(),
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(p)
}
