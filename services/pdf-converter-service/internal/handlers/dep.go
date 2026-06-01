package handlers

import (
	"pdf-converter-service/internal/config"
	"pdf-converter-service/internal/storage"
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
