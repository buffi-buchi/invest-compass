package profile

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
		logger:  logger.Named("api.v1.profile"),
	}
}

func (i *Implementation) Register(mux *chi.Mux) {
	mux.Get("/v1/profiles", i.GetProfiles)
}
