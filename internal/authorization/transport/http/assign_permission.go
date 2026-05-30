package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	authorizationapp "github.com/krishnaditya65/auth-server/internal/authorization/app"

	authctx "github.com/krishnaditya65/auth-server/internal/shared/context"
)

func (h *Handler) AssignPermission(
	w http.ResponseWriter,
	r *http.Request,
) {

	roleID := chi.URLParam(
		r,
		"roleID",
	)

	var req AssignPermissionRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(
		&req,
	)

	if err != nil {
		http.Error(
			w,
			"invalid request",
			http.StatusBadRequest,
		)
		return
	}

	p := authctx.MustPrincipal(
		r.Context(),
	)

	err = h.assignPermissionUseCase.Execute(
		r.Context(),
		authorizationapp.AssignPermissionToRoleInput{
			TenantID:       p.TenantID,
			RoleID:         roleID,
			PermissionName: req.Permission,
		},
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(
		http.StatusNoContent,
	)
}
