package routes

import (
	"net/http"

	"starter-kit-restapi-gonethttp/config"
	"starter-kit-restapi-gonethttp/internal/handlers"
	"starter-kit-restapi-gonethttp/internal/middleware"
)

func RegisterRoutes(cfg *config.Config, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler) http.Handler {
	mux := http.NewServeMux()
	healthHandler := handlers.NewHealthHandler()
	authMiddleware := middleware.Auth(cfg)
	rateLimit := middleware.RateLimit

	// Health
	mux.HandleFunc("GET /v1/health", healthHandler.HealthCheck)

	// Auth
	mux.HandleFunc("POST /v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /v1/auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("POST /v1/auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc("POST /v1/auth/verify-email", authHandler.VerifyEmail)
	mux.Handle("POST /v1/auth/send-verification-email", authMiddleware(http.HandlerFunc(authHandler.SendVerificationEmail)))

	// Users
	mux.Handle("POST /v1/users", authMiddleware(http.HandlerFunc(userHandler.CreateUser)))
	mux.Handle("GET /v1/users", authMiddleware(http.HandlerFunc(userHandler.GetUsers))) // New: Get List
	mux.Handle("GET /v1/users/{id}", authMiddleware(http.HandlerFunc(userHandler.GetUser)))
	mux.Handle("PATCH /v1/users/{id}", authMiddleware(http.HandlerFunc(userHandler.UpdateUser))) // New: Update
	mux.Handle("DELETE /v1/users/{id}", authMiddleware(http.HandlerFunc(userHandler.DeleteUser)))

	handler := middleware.Logger(mux)
	if cfg.Env == "production" {
		handler = rateLimit(handler)
	}

	return handler
}