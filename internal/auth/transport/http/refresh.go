package http

import (
	"encoding/json"
	"net/http"

	authapp "github.com/krishnaditya65/auth-server/internal/auth/app"
	sharederrors "github.com/krishnaditya65/auth-server/internal/shared/errors"
)

func (h *Handler) Refresh(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req RefreshRequest

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

	result, err :=
		h.refreshUseCase.Execute(
			r.Context(),
			authapp.RefreshInput{
				RefreshToken: req.RefreshToken,
			},
		)

	if err != nil {

		if err == sharederrors.ErrUnauthorized {
			http.Error(
				w,
				"unauthorized",
				http.StatusUnauthorized,
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

	resp := RefreshResponse{
		SessionID:    result.SessionID,
		TenantID:     result.TenantID,
		UserID:       result.UserID,
		RefreshToken: result.RefreshToken,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(
		resp,
	)
}
