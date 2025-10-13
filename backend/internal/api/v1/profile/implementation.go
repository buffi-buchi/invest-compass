package profile

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/buffi-buchi/invest-compass/backend/internal/api/middleware"
)

type Implementation struct {
	Unimplemented

	service Service
	auth    middleware.Middleware
	logger  *zap.Logger
}

func NewImplementation(service Service, auth middleware.Middleware, logger *zap.Logger) *Implementation {
	return &Implementation{
		service: service,
		auth:    auth,
		logger:  logger.Named("api.v1.profile"),
	}
}

func (i *Implementation) Register(mux *chi.Mux) {
	mux.Group(func(r chi.Router) {
		r.Use(i.auth)
		r.Get("/v1/profiles", i.GetProfiles)
	})
}
