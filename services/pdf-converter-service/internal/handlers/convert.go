package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"pdf-converter-service/internal/logger"
	"pdf-converter-service/internal/models"
	"pdf-converter-service/internal/services"
	"pdf-converter-service/internal/utils"

	"github.com/go-chi/chi/middleware"
)

func (h *Handler) ConvertHandler(log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Warn("missing file in request", map[string]any{
				"request_id": requestID,
				"error":      err.Error(),
			})
			return
		}
		defer file.Close()

		log.Info("file received", map[string]any{
			"request_id": requestID,
			"filename":   header.Filename,
			"size_bytes": header.Size,
		})

		if err := h.Storage.EnsureDir(); err != nil {
			log.Error("failed to create download folder", map[string]any{
				"request_id": requestID,
				"error":      err.Error(),
			})
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

		log.Info("ttl configured", map[string]any{
			"request_id":  requestID,
			"ttl_minutes": ttlMinutes,
		})

		timestamp := time.Now().UnixNano()
		safeFilename := utils.SanitizeFilename(header.Filename)
		filename := fmt.Sprintf("%d_%s", timestamp, safeFilename)
		pdfPath, err := h.Storage.SaveUploadedFile(file, filename)
		if err != nil {
			log.Error("failed to save pdf", map[string]any{
				"request_id": requestID,
				"filename":   filename,
				"error":      err.Error(),
			})

			return
		}

		log.Info("pdf saved", map[string]any{
			"request_id": requestID,
			"path":       pdfPath,
		})

		imgFiles, err := services.ConvertPDFToImages(
			pdfPath,
			h.Config.DownloadsDir,
			ttlMinutes,
		)
		if err != nil {
			log.Error("pdf conversion failed", map[string]any{
				"request_id": requestID,
				"path":       pdfPath,
				"error":      err.Error(),
			})
			http.Error(w, "conversion failed", http.StatusInternalServerError)
			return
		}

		log.Info("pdf converted to images", map[string]any{
			"request_id":  requestID,
			"path":        pdfPath,
			"image_count": len(imgFiles),
		})
		responses := []models.DownloadResponse{}
		for _, img := range imgFiles {
			info, err := h.Storage.StatFile(img)
			if err != nil {
				log.Error("failed to stat image", map[string]any{
					"request_id": requestID,
					"image":      img,
					"error":      err.Error(),
				})
				continue
			}

			url := fmt.Sprintf("/download/pdf/%s", filepath.Base(img))
			responses = append(responses, models.DownloadResponse{
				URL:       url,
				SizeBytes: info.Size(),
			})

			log.Info("download link generated", map[string]any{
				"request_id": requestID,
				"url":        url,
				"size_bytes": info.Size(),
			})
		}

		log.Info("request completed", map[string]any{
			"request_id":     requestID,
			"links_returned": len(responses),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responses)
	}
}
