package handlers

import "vaxen/api/internal/services"

// Handler holds all service dependencies for HTTP handlers.
// Individual handler files define methods on this struct.
type Handler struct {
	Services *services.Container
}

func NewHandler(svc *services.Container) *Handler {
	return &Handler{Services: svc}
}
