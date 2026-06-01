package handlers

import (
	"download-service/internal/config"
	"download-service/internal/storage"
)

type Handler struct {
	Config  config.Config
	Storage *storage.Storage
}

func New(cfg config.Config, storage *storage.Storage) *Handler {
	return &Handler{
		Config:  cfg,
		Storage: storage,
	}
}
