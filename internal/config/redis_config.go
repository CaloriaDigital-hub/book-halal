package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PGUrl         string
	RedisAddr     string
	RedisPassword string
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPass      string
	SMTPFrom      string
	HTTPPort      string
	StaticRoot    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to system environment variables")
	}

	return &Config{
		PGUrl:         getEnv("PG_URL", ""),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		SMTPHost:      getEnv("SMTP_HOST", ""),
		SMTPPort:      getEnv("SMTP_PORT", ""),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPass:      getEnv("SMTP_PASS", ""),
		SMTPFrom:      getEnv("SMTP_FROM", ""),
		HTTPPort:      getEnv("HTTP_PORT", ":8080"),
		StaticRoot:    getEnv("STATIC_ROOT", "./static"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}