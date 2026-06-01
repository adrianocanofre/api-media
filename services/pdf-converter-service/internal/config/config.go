package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerName      string
	ServerPort      string
	DownloadsDir    string
	APIGatewayURL   string
	CleanupInterval time.Duration
	DefaultTTL      int
	MaxTTL          int
	Debug           bool
}

func LoadConfig() Config {
	c := Config{}

	if debug, err := strconv.ParseBool(getEnv("DEBUG", "false")); err == nil {
		c.Debug = debug
	}

	c.ServerName = getEnv("SERVER_NAME", "pdf-service")
	c.ServerPort = getEnv("SERVER_PORT", "8081")

	c.DownloadsDir = getEnv("PDF_DOWNLOAD_DIR", "/tmp/downloads/")
	c.APIGatewayURL = getEnv("API_GATEWAY_URL", "localhost:8080")

	intervalStr := getEnv("PDF_CLEANUP_INTERVAL", "1m")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		interval = time.Minute
	}
	c.CleanupInterval = interval

	if ttl, err := strconv.Atoi(getEnv("PDF_DEFAULT_TTL_MINUTES", "5")); err == nil {
		c.DefaultTTL = ttl
	} else {
		c.DefaultTTL = 5
	}

	if maxTTL, err := strconv.Atoi(getEnv("PDF_MAX_TTL_MINUTES", "60")); err == nil {
		c.MaxTTL = maxTTL
	} else {
		c.MaxTTL = 60
	}

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
	log.Printf("API Gateway URL: %s", cfg.APIGatewayURL)
	log.Printf("Cleanup Interval: %s", cfg.CleanupInterval)
	log.Printf("Default TTL: %d", cfg.DefaultTTL)
	log.Printf("Max TTL: %d", cfg.MaxTTL)

}
