package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Env      string
	Database DatabaseConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret                       string
	AccessExpirationMinutes      int
	RefreshExpirationDays        int
	ResetPasswordExpirationMinutes int
	VerifyEmailExpirationMinutes   int
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// LoadConfig loads environment variables from .env file
func LoadConfig() *Config {
	// Load .env file if it exists (ignore error if not found, useful for Docker/Prod)
	_ = godotenv.Load()

	return &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("GO_ENV", "development"),
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "starter_kit_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:                       getEnv("JWT_SECRET", "secret"),
			AccessExpirationMinutes:      getEnvAsInt("JWT_ACCESS_EXPIRATION_MINUTES", 30),
			RefreshExpirationDays:        getEnvAsInt("JWT_REFRESH_EXPIRATION_DAYS", 30),
			ResetPasswordExpirationMinutes: getEnvAsInt("JWT_RESET_PASSWORD_EXPIRATION_MINUTES", 10),
			VerifyEmailExpirationMinutes:   getEnvAsInt("JWT_VERIFY_EMAIL_EXPIRATION_MINUTES", 10),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getEnvAsInt("SMTP_PORT", 587),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("EMAIL_FROM", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}