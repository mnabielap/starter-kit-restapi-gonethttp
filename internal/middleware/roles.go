package middleware

import (
	"net/http"

	"starter-kit-restapi-gonethttp/internal/services"
	"starter-kit-restapi-gonethttp/pkg/response"

	"github.com/google/uuid"
)

// RequireAdmin ensures the authenticated user has the 'admin' role.
func RequireAdmin(service services.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get UserID from context (set by Auth middleware)
			userIDStr, ok := r.Context().Value(UserIDKey).(string)
			if !ok {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// 2. Fetch User from DB
			id, err := uuid.Parse(userIDStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Invalid User ID")
				return
			}

			user, err := service.GetUserByID(id)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "User not found")
				return
			}

			// 3. Check Role
			if user.Role != "admin" {
				response.Error(w, http.StatusForbidden, "Forbidden: Admins only")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdminOrSelf ensures the user is admin OR accessing their own resource (for Get One).
func RequireAdminOrSelf(service services.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDStr, ok := r.Context().Value(UserIDKey).(string)
			if !ok {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			targetID := r.PathValue("id")

			// Check if Self
			if targetID == userIDStr {
				next.ServeHTTP(w, r)
				return
			}

			// If not self, Check if Admin
			id, _ := uuid.Parse(userIDStr)
			user, err := service.GetUserByID(id)
			if err != nil || user.Role != "admin" {
				response.Error(w, http.StatusForbidden, "Forbidden: Access denied")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}