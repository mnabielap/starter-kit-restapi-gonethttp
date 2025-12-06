package middleware

import (
	"context"
	"net/http"
	"strings"

	"starter-kit-restapi-gonethttp/config"
	"starter-kit-restapi-gonethttp/pkg/response"
	"starter-kit-restapi-gonethttp/pkg/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "Missing Authorization header")
				return
			}

			// Format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(w, http.StatusUnauthorized, "Invalid Authorization header format")
				return
			}

			tokenString := parts[1]
			claims, err := utils.ValidateToken(tokenString, cfg.JWT.Secret)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			if claims.Type != "access" {
				response.Error(w, http.StatusUnauthorized, "Invalid token type")
				return
			}

			// Add UserID to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.Sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}