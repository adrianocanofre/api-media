package services

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"pdf-converter-service/internal/models"
	"pdf-converter-service/internal/storage"
)

func ConvertPDFToImages(pdfPath, outputDir string, ttlMinutes int) ([]string, error) {
	baseName := filepath.Base(pdfPath)
	baseName = baseName[:len(baseName)-len(filepath.Ext(baseName))]

	cmd := exec.Command("pdftoppm", "-png", "-r", "200", pdfPath, filepath.Join(outputDir, baseName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(outputDir, baseName+"-*.png"))
	if err != nil {
		return nil, err
	}

	metadata := []models.ImageMetadata{}
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			log.Printf("[WARN] Failed to stat image %s: %v", f, err)
			continue
		}

		metadata = append(metadata, models.ImageMetadata{
			Filename:   filepath.Base(f),
			PDFFile:    filepath.Base(pdfPath),
			CreatedAt:  info.ModTime(),
			TTLMinutes: ttlMinutes,
		})
	}

	if err := storage.AppendMetadataToFile(filepath.Join(outputDir, "images-metadata.json"), metadata); err != nil {
		log.Printf("[ERROR] Failed to write metadata JSON: %v", err)
	}

	return files, nil
}
