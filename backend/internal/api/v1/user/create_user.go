package user

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/buffi-buchi/invest-compass/backend/internal/api"
	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

func (i *Implementation) CreateUser(w http.ResponseWriter, r *http.Request) {
	request, err := api.DecodeRequest[CreateUserRequest](r)
	if err != nil {
		i.logger.Error("Decode create user request", zap.Error(err))

		api.EncodeErrorf(w, http.StatusBadRequest, "Invalid request: %s", err)

		return
	}

	user, err := i.service.Create(r.Context(), model.User{
		Email:    request.Email,
		Password: request.Password,
	})
	if errors.Is(err, model.ErrAlreadyExists) {
		i.logger.Error("User already exists", zap.Error(err))

		api.EncodeErrorf(w, http.StatusConflict, "User already exists")

		return
	}
	if err != nil {
		i.logger.Error("Create user", zap.Error(err))

		api.EncodeErrorf(w, http.StatusInternalServerError, "Create user error")

		return
	}

	response := CreateUserResponse{
		Id:         user.ID,
		Email:      user.Email,
		CreateTime: user.CreateTime,
	}

	api.EncodeSuccess(w, http.StatusCreated, response)
}
