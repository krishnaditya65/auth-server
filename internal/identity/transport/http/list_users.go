package http

import (
	"encoding/json"
	"net/http"

	authctx "github.com/krishnaditya65/auth-server/internal/shared/context"
)

func (h *Handler) ListUsers(
	w http.ResponseWriter,
	r *http.Request,
) {

	p := authctx.MustPrincipal(
		r.Context(),
	)

	users, err := h.listUsersUseCase.Execute(
		r.Context(),
		p.TenantID,
	)

	if err != nil {
		http.Error(
			w,
			"failed to list users",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(
		users,
	)
}
