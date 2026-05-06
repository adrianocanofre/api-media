package storage

import (
	"encoding/json"
	"os"

	"pdf-converter-service/internal/models"
)

var metadataMutex = make(chan struct{}, 1) // simples lock

func AppendMetadataToFile(path string, newData []models.ImageMetadata) error {
	metadataMutex <- struct{}{}
	defer func() { <-metadataMutex }()

	var existing []models.ImageMetadata

	if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
		json.Unmarshal(data, &existing)
	}

	existing = append(existing, newData...)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(existing)
}
