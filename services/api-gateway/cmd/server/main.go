package main

import (
	"net/http"
	"time"

	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/logger"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.LoadConfig()
	config.PrintStartupConfig(cfg)

	log := logger.New(cfg.ServerName)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(handlers.LoggerMiddleware(log))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(`{
		"service":"api-gateway",
		"status":"running"
	}`))
	})
	r.Get("/health", handlers.HealthHandler)
	r.Post("/pdf/pdf-to-image", handlers.PdfProxyHandler(log, cfg))
	r.Get("/download/pdf/{filename}", handlers.DownloadProxyHandler(log, cfg))

	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Error("server failed", map[string]any{
			"error": err.Error(),
		})
	}
}
