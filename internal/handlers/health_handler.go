package handlers

import (
	"net/http"

	"starter-kit-restapi-gonethttp/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"message": "Server is running",
	})
}