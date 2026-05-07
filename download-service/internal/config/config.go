package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerName   string
	ServerPort   string
	DownloadsDir string
	Debug        bool
}

func LoadConfig() Config {
	c := Config{}

	if debug, err := strconv.ParseBool(getEnv("DEBUG", "false")); err == nil {
		c.Debug = debug
	}

	c.ServerName = getEnv("SERVER_NAME", "download-service")
	c.ServerPort = getEnv("SERVER_PORT", "8082")
	c.DownloadsDir = getEnv("PDF_DOWNLOAD_DIR", "/tmp/downloads/")

	return c
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func PrintStartupConfig(cfg Config) {
	log.Printf("Server Name: %s", cfg.ServerName)
	log.Printf("Server Port: %s", cfg.ServerPort)
	log.Printf("Downloads Dir: %s", cfg.DownloadsDir)

	if !cfg.Debug {
		return
	}

	log.Println("Debug mode enabled")
}
