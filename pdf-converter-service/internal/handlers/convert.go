package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"pdf-converter-service/internal/models"
	"pdf-converter-service/internal/services"
	"pdf-converter-service/internal/utils"
)

func (h *Handler) ConvertHandler(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := h.Storage.EnsureDir(); err != nil {
		http.Error(w, "failed to create download folder", http.StatusInternalServerError)
		return
	}

	// TTL default 5 min
	ttlMinutes := h.Config.DefaultTTL
	if v := r.FormValue("ttl_minutes"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			if parsed > h.Config.MaxTTL {
				ttlMinutes = 60
			} else if parsed > 0 {
				ttlMinutes = parsed
			}
		}
	}

	log.Printf("[INFO] PDF TTL configured ttl_minutes=%d", ttlMinutes)

	timestamp := time.Now().UnixNano()
	safeFilename := utils.SanitizeFilename(header.Filename)
	filename := fmt.Sprintf("%d_%s", timestamp, safeFilename)
	pdfPath, err := h.Storage.SaveUploadedFile(file, filename)
	if err != nil {
		http.Error(w, "failed to save PDF", http.StatusInternalServerError)
		return
	} else {
		log.Printf("[INFO] PDF  saved in %s", pdfPath)
	}

	imgFiles, err := services.ConvertPDFToImages(
		pdfPath,
		h.Config.DownloadsDir,
		ttlMinutes,
	)
	if err != nil {
		log.Printf("PDF conversion failed for %s: %v", pdfPath, err)
		http.Error(w, "conversion failed", http.StatusInternalServerError)
		return
	}

	responses := []models.DownloadResponse{}
	for _, img := range imgFiles {
		info, err := h.Storage.StatFile(img)
		if err != nil {
			log.Printf("failed to stat image %s: %v", img, err)
			continue
		}

		url := fmt.Sprintf("http://%s/download/pdf/%s", h.Config.APIGatewayURL, filepath.Base(img))
		responses = append(responses, models.DownloadResponse{
			URL:       url,
			SizeBytes: info.Size(),
		})

		log.Printf("download link generated url=%s size_bytes=%d", url, info.Size())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
}
