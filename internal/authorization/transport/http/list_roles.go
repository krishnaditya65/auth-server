package http

import (
	"encoding/json"
	"net/http"

	authctx "github.com/krishnaditya65/auth-server/internal/shared/context"
)

func (h *Handler) ListRoles(
	w http.ResponseWriter,
	r *http.Request,
) {

	p := authctx.MustPrincipal(
		r.Context(),
	)

	roles, err := h.listRolesUseCase.Execute(
		r.Context(),
		p.TenantID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(
		roles,
	)
}
