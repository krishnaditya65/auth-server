package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	authctx "github.com/krishnaditya65/auth-server/internal/shared/context"
)

func (h *Handler) ListRolePermissions(
	w http.ResponseWriter,
	r *http.Request,
) {

	roleID := chi.URLParam(
		r,
		"roleID",
	)

	p := authctx.MustPrincipal(
		r.Context(),
	)

	permissions, err :=
		h.listRolePermissionsUseCase.Execute(
			r.Context(),
			p.TenantID,
			roleID,
		)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(
		w,
	).Encode(
		permissions,
	)
}
