package main

import (
	"fmt"
	"net/http"
	"os"

	"starter-kit-restapi-gonethttp/config"
	"starter-kit-restapi-gonethttp/internal/handlers"
	"starter-kit-restapi-gonethttp/internal/repository"
	"starter-kit-restapi-gonethttp/internal/routes"
	"starter-kit-restapi-gonethttp/internal/services"
	"starter-kit-restapi-gonethttp/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.Env)
	logger.Log.Info("Starting server...", "env", cfg.Env)

	config.ConnectDB(cfg)

	userRepo := repository.NewUserRepository(config.DB)
	tokenRepo := repository.NewTokenRepository(config.DB)

	tokenService := services.NewTokenService(tokenRepo, cfg)
	emailService := services.NewEmailService(cfg)
	userService := services.NewUserService(userRepo)
	
	authService := services.NewAuthService(userRepo, tokenRepo, tokenService, emailService, cfg)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	router := routes.RegisterRoutes(cfg, authHandler, userHandler)

	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	logger.Log.Info("Server listening", "address", serverAddr)
	
	err := http.ListenAndServe(serverAddr, router)
	if err != nil {
		logger.Log.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}