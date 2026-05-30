package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	authctx "github.com/krishnaditya65/auth-server/internal/shared/context"

	identityapp "github.com/krishnaditya65/auth-server/internal/identity/app"
)

func (h *Handler) GetUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := chi.URLParam(
		r,
		"userID",
	)

	p := authctx.MustPrincipal(
		r.Context(),
	)

	user, err := h.getUserUseCase.Execute(
		r.Context(),
		identityapp.GetUserInput{
			TenantID: p.TenantID,
			UserID:   userID,
		},
	)

	if err != nil {
		http.Error(
			w,
			"user not found",
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(
		user,
	)
}
