package handlers

import (
	"log"
	"net/http"
	"path/filepath"
)

func (h *Handler) DownloadPDFHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/download/pdf/"):]

	path := filepath.Join(h.Storage.DownloadsDir, id)
	log.Printf(path)
	if _, err := h.Storage.StatFile(path); err != nil {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, path)
}
