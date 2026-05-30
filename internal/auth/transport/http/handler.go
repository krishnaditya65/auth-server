package http

import (
	"encoding/json"
	"net/http"

	authapp "github.com/krishnaditya65/auth-server/internal/auth/app"
	sharederrors "github.com/krishnaditya65/auth-server/internal/shared/errors"
)

type Handler struct {
	registerUseCase *authapp.RegisterUseCase
	loginUseCase    *authapp.LoginUseCase
	refreshUseCase  *authapp.RefreshUseCase
}

func NewHandler(
	registerUseCase *authapp.RegisterUseCase,
	loginUseCase *authapp.LoginUseCase,
	refreshUseCase *authapp.RefreshUseCase,
) *Handler {

	return &Handler{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
		refreshUseCase:  refreshUseCase,
	}
}

func (h *Handler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	err = h.registerUseCase.Execute(
		r.Context(),
		authapp.RegisterInput{
			Email:    req.Email,
			Password: req.Password,
		},
	)

	if err != nil {
		if err == sharederrors.ErrConflict {
			http.Error(w, "email already exists", http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{
		Message: "registration successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(resp)
}
