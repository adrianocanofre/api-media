package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		recorder.Header().Set("X-Request-Id", requestID)

		next.ServeHTTP(recorder, r)

		durationMs := time.Since(start).Milliseconds()
		log.Printf("method=%s path=%s status=%d total_ms=%d request_id=%s remote_addr=%s",
			r.Method, r.URL.Path, recorder.status, durationMs, requestID, r.RemoteAddr)
	})
}
