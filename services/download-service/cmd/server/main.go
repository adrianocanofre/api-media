package main

import (
	"net/http"
	"time"

	"download-service/internal/config"
	"download-service/internal/handlers"
	"download-service/internal/logger"
	"download-service/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.LoadConfig()
	config.PrintStartupConfig(cfg)

	st := storage.New(cfg.DownloadsDir)
	h := handlers.New(cfg, st)

	log := logger.New(cfg.ServerName)
	// Setup routes
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(handlers.LoggerMiddleware(log))

	r.Get("/health", handlers.HealthHandler)
	r.Get("/download/pdf/{filename}", h.DownloadPDFHandler)

	// Wrap with logging middleware

	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Error("server failed", map[string]any{
			"error": err.Error(),
		})
	}
}
