package handlers

import (
	"net/http"
	"time"

	"download-service/internal/logger"

	"github.com/go-chi/chi/v5/middleware"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := middleware.GetReqID(r.Context())

			log.Info("request started", map[string]any{
				"method":     r.Method,
				"path":       r.URL.Path,
				"request_id": requestID,
				"ip":         r.RemoteAddr,
			})

			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)

			log.Info("request completed", map[string]any{
				"method":     r.Method,
				"path":       r.URL.Path,
				"request_id": requestID,
				"status":     rec.status,
				"duration":   time.Since(start).String(),
			})
		})
	}
}
