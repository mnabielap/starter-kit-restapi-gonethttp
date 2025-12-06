package handlers

import (
	"encoding/json"
	"net/http"

	"starter-kit-restapi-gonethttp/internal/middleware"
	"starter-kit-restapi-gonethttp/internal/models"
	"starter-kit-restapi-gonethttp/internal/services"
	"starter-kit-restapi-gonethttp/pkg/response"
	"starter-kit-restapi-gonethttp/pkg/utils"

	"github.com/google/uuid"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req services.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := utils.ValidateStruct(req); errs != nil {
		response.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "Validation error", "errors": errs})
		return
	}

	user, tokens, err := h.service.Register(req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	user, tokens, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.service.Logout(req.RefreshToken); err != nil {
		response.Error(w, http.StatusNotFound, "Not found")
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := h.service.ForgotPassword(req.Email); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		response.Error(w, http.StatusBadRequest, "Token is required")
		return
	}

	var req struct {
		Password string `json:"password" validate:"required,min=8"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := h.service.ResetPassword(token, req.Password); err != nil {
		response.Error(w, http.StatusUnauthorized, "Password reset failed")
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	id, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Create a partial user object to pass to service
	// In a real app, you might want to fetch the full user from DB if Email is needed and not in context
	user := &models.User{ID: id}

	if err := h.service.SendVerificationEmail(user); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		response.Error(w, http.StatusBadRequest, "Token is required")
		return
	}

	if err := h.service.VerifyEmail(token); err != nil {
		response.Error(w, http.StatusUnauthorized, "Email verification failed")
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}