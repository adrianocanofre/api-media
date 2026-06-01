package models

import (
	"time"
)

type DownloadResponse struct {
	URL       string `json:"url"`
	SizeBytes int64  `json:"size_bytes"`
}

type ImageMetadata struct {
	Filename   string    `json:"filename"`
	CreatedAt  time.Time `json:"created_at"`
	PDFFile    string    `json:"pdf_file"`
	TTLMinutes int       `json:"ttl_minutes"`
}
