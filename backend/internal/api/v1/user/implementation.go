package user

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Implementation struct {
	Unimplemented

	service Service
	logger  *zap.Logger
}

func NewImplementation(service Service, logger *zap.Logger) *Implementation {
	return &Implementation{
		service: service,
		logger:  logger.Named("api.v1.user"),
	}
}

func (i *Implementation) Register(mux *chi.Mux) {
	mux.Post("/v1/users", i.CreateUser)
}
