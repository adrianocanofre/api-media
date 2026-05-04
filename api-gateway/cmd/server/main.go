package main

import (
	"net/http"
	"time"

	"api-gateway/internal/config"
	"api-gateway/internal/httpserver"
	"api-gateway/internal/logger"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.LoadConfig()

	log := logger.New(cfg.ServerName)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(httpserver.LoggerMiddleware(log))

	r.Get("/healthz", httpserver.HealthHandler)
	r.Post("/pdf/pdf-to-image", httpserver.PdfProxyHandler(log, cfg))
	r.Get("/download/pdf/{filename}", httpserver.DownloadProxyHandler(log, cfg))

	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Error("server failed", map[string]any{
			"error": err.Error(),
		})
	}
}
