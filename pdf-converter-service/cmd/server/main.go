package main

import (
	"net/http"
	"time"

	"pdf-converter-service/internal/config"
	"pdf-converter-service/internal/handlers"
	"pdf-converter-service/internal/logger"
	"pdf-converter-service/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.LoadConfig()
	config.PrintStartupConfig(cfg)

	storage.StartImageCleanupJob(cfg.CleanupInterval, cfg.DownloadsDir)

	st := storage.New(cfg.DownloadsDir)
	h := handlers.New(cfg, st)

	log := logger.New(cfg.ServerName)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(handlers.LoggerMiddleware(log))

	r.Get("/health", handlers.HealthHandler)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/convert/pdf-to-image", h.ConvertHandler(log))
	})

	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Error("server failed", map[string]any{
			"error": err.Error(),
		})
	}
}
