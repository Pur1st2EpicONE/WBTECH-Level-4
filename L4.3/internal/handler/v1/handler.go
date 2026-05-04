package v1

import (
	"L4.3/internal/service"
	"L4.3/pkg/logger"
)

type Handler struct {
	service service.Service
	logger  logger.Logger
}

func NewHandler(service service.Service, logger logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}
