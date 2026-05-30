package http

import (
	"encoding/json"
	"net/http"

	identityapp "github.com/krishnaditya65/auth-server/internal/identity/app"

	authctx "github.com/krishnaditya65/auth-server/internal/shared/context"
)

type CreateUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	RoleName    string `json:"role_name"`
}

func (h *Handler) CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req CreateUserRequest

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

	err = h.createUserUseCase.Execute(
		r.Context(),
		identityapp.CreateUserInput{
			TenantID:    p.TenantID,
			Email:       req.Email,
			Password:    req.Password,
			DisplayName: req.DisplayName,
			RoleName:    req.RoleName,
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
		http.StatusCreated,
	)
}
