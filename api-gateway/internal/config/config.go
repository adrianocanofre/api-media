package config

import "os"

type Config struct {
	ServerName      string
	ServerPort      string
	PdfService      string
	DownloadService string
}

func LoadConfig() *Config {
	return &Config{
		ServerName:      getEnv("SERVER_NAME", "api-gateway"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		PdfService:      getEnv("PDF_SERVICE_URL", "http://localhost:8081"),
		DownloadService: getEnv("DOWNLOAD_SERVICE_URL", "http://localhost:8082"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
