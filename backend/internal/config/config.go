package config

import (
	"log"
	"os"
	"strconv"
	"strings"

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
	RAGServiceURL string
	BaseURL       string // Public URL of this server, e.g. http://localhost:8090
	PrivateRoot   string // Приватные файлы (чеки об оплате) — не раздаются через /static

	// Telegram-бот для проверки оплат. Пустой токен — бот выключен.
	TelegramBotToken         string
	TelegramChatID           int64   // группа с темами (supergroup), например -1001234567890
	TelegramPaymentsThreadID int64   // тема «Платежи»
	TelegramHistoryThreadID  int64   // тема «История»
	TelegramAdminIDs         []int64 // кто может нажимать кнопки; пусто — любой участник группы
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
		RAGServiceURL: getEnv("RAG_SERVICE_URL", "http://localhost:8001"),
		BaseURL:       getEnv("BASE_URL", "http://localhost:8090"),
		PrivateRoot:   getEnv("PRIVATE_ROOT", "./storage"),

		TelegramBotToken:         getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:           getEnvInt64("TELEGRAM_CHAT_ID"),
		TelegramPaymentsThreadID: getEnvInt64("TELEGRAM_PAYMENTS_THREAD_ID"),
		TelegramHistoryThreadID:  getEnvInt64("TELEGRAM_HISTORY_THREAD_ID"),
		TelegramAdminIDs:         getEnvInt64List("TELEGRAM_ADMIN_IDS"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
func getEnvInt64(key string) int64 {
	raw := strings.TrimSpace(getEnv(key, ""))
	if raw == "" {
		return 0
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		log.Fatalf("config: %s must be an integer, got %q", key, raw)
	}
	return v
}

// getEnvInt64List: "123, 456" → [123 456]
func getEnvInt64List(key string) []int64 {
	var result []int64
	for _, part := range strings.Split(getEnv(key, ""), ",") {
		if part = strings.TrimSpace(part); part == "" {
			continue
		}
		v, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			log.Fatalf("config: %s must be comma-separated integers, got %q", key, part)
		}
		result = append(result, v)
	}
	return result
}
