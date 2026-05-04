package main

import (
	"net/http"
	"time"

	"api-gateway/internal/httpserver"
	"api-gateway/internal/logger"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	log := logger.New("api-gateway")

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(httpserver.LoggerMiddleware(log))

	r.Get("/healthz", httpserver.HealthHandler)
	r.Post("/pdf/pdf-to-image", httpserver.PdfProxyHandler(log))
	r.Get("/download/pdf/{filename}", httpserver.DownloadProxyHandler(log))

	log.Info("server started", map[string]any{
		"addr": ":8080",
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Error("server failed", map[string]any{
			"error": err.Error(),
		})
	}
}
