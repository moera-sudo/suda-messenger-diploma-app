package config

import (
	"os"
	"strconv"
)

type EmailConfig struct {
	// SMTP Settings
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string
}

func Load() *EmailConfig {
	return &EmailConfig{
		SMTPHost: getEnv("SMTP_HOST", ""),
		SMTPPort: getEnvInt("SMTP_PORT", 2525),
		SMTPUser: getEnv("SMTP_USER", ""),
		SMTPPass: getEnv("SMTP_PASS", ""),
		SMTPFrom: getEnv("SMTP_FROM", "suda@messenger.com"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}