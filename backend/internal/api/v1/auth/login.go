package auth

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/buffi-buchi/invest-compass/backend/internal/api"
)

func (i *Implementation) Login(w http.ResponseWriter, r *http.Request) {
	request, err := api.DecodeRequest[LoginRequest](r)
	if err != nil {
		i.logger.Error("Decode login request", zap.Error(err))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request: %s", err)

		return
	}

	token, err := i.service.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		i.logger.Error("Login", zap.Error(err))

		api.EncodeErrorf(w, http.StatusUnauthorized, "Unauthenticated")

		return
	}

	response := LoginResponse{
		Token: token,
	}

	api.EncodeSuccess(w, http.StatusOK, response)
}
