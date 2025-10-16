package user

import "go.uber.org/zap"

type Implementation struct {
	Unimplemented

	service Service
	logger  *zap.SugaredLogger
}

func NewImplementation(service Service, logger *zap.SugaredLogger) *Implementation {
	return &Implementation{
		service: service,
		logger:  logger.Named("api.v1.user"),
	}
}
