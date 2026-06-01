package storage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"pdf-converter-service/internal/models"
)

func CleanupExpiredImages(metadataPath, imagesDir string) {
	metadataMutex <- struct{}{}
	defer func() { <-metadataMutex }()

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Printf("[ERROR] Failed to read metadata file: %v", err)
		return
	}

	var metadata []models.ImageMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		log.Printf("[ERROR] Failed to parse metadata JSON: %v", err)
		return
	}

	now := time.Now()
	kept := []models.ImageMetadata{}

	for _, item := range metadata {
		itemTTL := time.Duration(item.TTLMinutes) * time.Minute
		age := now.Sub(item.CreatedAt)

		if age <= itemTTL {
			kept = append(kept, item)
			continue
		}

		// delete image
		imagePath := filepath.Join(imagesDir, item.Filename)
		os.Remove(imagePath)

		// delete PDF if last image expired
		pdfPath := filepath.Join(imagesDir, item.PDFFile)
		os.Remove(pdfPath)

		log.Printf("[INFO] Expired image and PDF removed filename=%s pdf=%s", item.Filename, item.PDFFile)
	}

	file, err := os.Create(metadataPath)
	if err != nil {
		log.Printf("[ERROR] Failed to rewrite metadata file: %v", err)
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(kept)
}

func StartImageCleanupJob(interval time.Duration, downloadsDir string) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			path := filepath.Join(downloadsDir, "images-metadata.json")
			log.Printf("[INFO] Running cleanup job metadata=%s", path)

			CleanupExpiredImages(path, downloadsDir)
		}
	}()
}
