package http

import (
	"encoding/json"
	"net/http"

	authorizationapp "github.com/krishnaditya65/auth-server/internal/authorization/app"

	authctx "github.com/krishnaditya65/auth-server/internal/shared/context"
	sharederrors "github.com/krishnaditya65/auth-server/internal/shared/errors"
)

func (h *Handler) CreateRole(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req CreateRoleRequest

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

	err = h.createRoleUseCase.Execute(
		r.Context(),
		authorizationapp.CreateRoleInput{
			TenantID:    p.TenantID,
			Name:        req.Name,
			Description: req.Description,
		},
	)

	if err != nil {

		if err == sharederrors.ErrConflict {
			http.Error(
				w,
				"role already exists",
				http.StatusConflict,
			)
			return
		}

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(
		http.StatusCreated,
	)
}
