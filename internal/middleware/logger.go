package middleware

import (
	"net/http"
	"time"

	"starter-kit-restapi-gonethttp/pkg/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		
		next.ServeHTTP(wrappedWriter, r)

		logger.Log.Info("Request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrappedWriter.status,
			"duration", time.Since(start).String(),
			"ip", r.RemoteAddr,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}