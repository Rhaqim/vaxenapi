package handlers

import (
	"vaxen/api/internal/config"
	"vaxen/api/internal/services"
)

// Handler holds all dependencies for HTTP handlers.
type Handler struct {
	Services *services.Container
	Config   *config.Config
}

func NewHandler(svc *services.Container, cfg *config.Config) *Handler {
	return &Handler{Services: svc, Config: cfg}
}
