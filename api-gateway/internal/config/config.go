package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerName      string
	ServerPort      string
	PdfService      string
	DownloadService string
	Debug           bool
}

func LoadConfig() *Config {
	c := &Config{}

	if debug, err := strconv.ParseBool(getEnv("DEBUG", "false")); err == nil {
		c.Debug = debug
	}

	c.ServerName = getEnv("SERVER_NAME", "api-gateway")
	c.ServerPort = getEnv("SERVER_PORT", "8080")
	c.PdfService = getEnv("PDF_SERVICE_URL", "http://localhost:8081")
	c.DownloadService = getEnv("DOWNLOAD_SERVICE_URL", "http://localhost:8082")

	return c
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func PrintStartupConfig(cfg Config) {
	log.Printf("Server Name: %s", cfg.ServerName)
	log.Printf("Server Port: %s", cfg.ServerPort)

	if !cfg.Debug {
		return
	}

	log.Println("Debug mode enabled")
	log.Printf("PDF Server: %s", cfg.PdfService)
	log.Printf("Download Server: %s", cfg.DownloadService)

}
