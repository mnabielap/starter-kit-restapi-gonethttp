package routes

import (
	"net/http"

	"starter-kit-restapi-gonethttp/config"
	"starter-kit-restapi-gonethttp/internal/handlers"
	"starter-kit-restapi-gonethttp/internal/middleware"
	"starter-kit-restapi-gonethttp/internal/services"
)

func RegisterRoutes(cfg *config.Config, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, userService services.UserService) http.Handler {
	mux := http.NewServeMux()
	healthHandler := handlers.NewHealthHandler()
	authMiddleware := middleware.Auth(cfg)
	rateLimit := middleware.RateLimit
	
	// Role Middleware
	requireAdmin := middleware.RequireAdmin(userService)
	requireAdminOrSelf := middleware.RequireAdminOrSelf(userService)

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

	// Users (Protected with RBAC)
	
	// Create User: Admin Only
	mux.Handle("POST /v1/users", authMiddleware(requireAdmin(http.HandlerFunc(userHandler.CreateUser))))
	
	// Get List: Admin Only
	mux.Handle("GET /v1/users", authMiddleware(requireAdmin(http.HandlerFunc(userHandler.GetUsers))))
	
	// Get One: Admin OR Self
	mux.Handle("GET /v1/users/{id}", authMiddleware(requireAdminOrSelf(http.HandlerFunc(userHandler.GetUser))))
	
	// Update: Admin Only (Strict CRUD)
	// If you want users to update themselves, use requireAdminOrSelf here.
	mux.Handle("PATCH /v1/users/{id}", authMiddleware(requireAdmin(http.HandlerFunc(userHandler.UpdateUser))))
	
	// Delete: Admin Only
	mux.Handle("DELETE /v1/users/{id}", authMiddleware(requireAdmin(http.HandlerFunc(userHandler.DeleteUser))))

	handler := middleware.Logger(mux)
	if cfg.Env == "production" {
		handler = rateLimit(handler)
	}

	return handler
}