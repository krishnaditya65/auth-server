package http

import (
	"encoding/json"
	"net/http"

	authapp "github.com/krishnaditya65/auth-server/internal/auth/app"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req LoginRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&req)

	if err != nil {
		http.Error(
			w,
			"invalid request",
			http.StatusBadRequest,
		)
		return
	}

	result, err := h.loginUseCase.Execute(
		r.Context(),
		authapp.LoginInput{
			Email:    req.Email,
			Password: req.Password,
		},
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

	json.NewEncoder(w).Encode(result)
}
