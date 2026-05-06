package main

import (
	"log"
	"net/http"

	"pdf-converter-service/internal/config"
	"pdf-converter-service/internal/handlers"
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

	r := chi.NewRouter()

	// Middlewares globais
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// Setup routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/convert/pdf-to-image", h.ConvertHandler)
	})

	r.Get("/health", handlers.HealthHandler)

	// Wrap with logging middleware

	log.Println("Servidor rodando em http://localhost:" + cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
