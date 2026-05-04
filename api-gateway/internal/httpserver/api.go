package httpserver

import (
	"io"
	"net/http"
	"time"

	"api-gateway/internal/logger"

	"github.com/go-chi/chi/v5/middleware"
)

const pdfService = "http://localhost:8082"
const downloadService = "http://localhost:8083"

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func PdfProxyHandler(log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())

		log.Info("proxying request to pdf service", map[string]any{
			"request_id": requestID,
			"target":     pdfService + "/convert/pdf-to-image",
		})
		req, err := http.NewRequestWithContext(
			r.Context(),
			http.MethodPost,
			pdfService+"/convert/pdf-to-image",
			r.Body,
		)
		if err != nil {
			log.Error("failed to create request", map[string]any{
				"path":  r.URL.Path,
				"error": err.Error(),
			})
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}

		for k, v := range r.Header {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}

		log.Info("forwarding headers to pdf service", map[string]any{
			"request_id":   requestID,
			"header_count": len(r.Header),
		})

		start := time.Now()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error("pdf service unavailable", map[string]any{
				"request_id": requestID,
				"path":       r.URL.Path,
				"target":     pdfService,
				"error":      err.Error(),
			})
			http.Error(w, "pdf service unavailable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		elapsed := time.Since(start)

		log.Info("pdf service responded", map[string]any{
			"request_id": requestID,
			"status":     resp.StatusCode,
			"duration":   elapsed.String(),
		})

		if resp.StatusCode >= 400 {
			log.Warn("pdf service returned error status", map[string]any{
				"request_id": requestID,
				"status":     resp.StatusCode,
			})
		}

		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func DownloadProxyHandler(log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		targetURL := downloadService + r.URL.Path

		log.Info("proxying request to download service", map[string]any{
			"request_id": requestID,
			"target":     targetURL,
		})

		req, err := http.NewRequestWithContext(
			r.Context(),
			http.MethodGet,
			targetURL,
			nil,
		)
		if err != nil {
			log.Error("failed to create request", map[string]any{
				"path":  r.URL.Path,
				"error": err.Error(),
			})
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}

		for k, v := range r.Header {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}

		log.Info("forwarding headers to download service", map[string]any{
			"request_id":   requestID,
			"header_count": len(r.Header),
		})

		start := time.Now()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error("download service unavailable", map[string]any{
				"request_id": requestID,
				"path":       r.URL.Path,
				"target":     downloadService,
				"error":      err.Error(),
			})
			http.Error(w, "download service unavailable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		elapsed := time.Since(start)

		log.Info("download service responded", map[string]any{
			"request_id": requestID,
			"status":     resp.StatusCode,
			"duration":   elapsed.String(),
		})

		if resp.StatusCode >= 400 {
			log.Warn("download service returned error status", map[string]any{
				"request_id": requestID,
				"status":     resp.StatusCode,
			})
		}

		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
