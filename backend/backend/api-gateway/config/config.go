package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	MessengerHost          string
	TxHost                 string
	MediaHost              string
	JWTSecret              string
	GatewaySignatureSecret string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {

		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("No .env file found, using system env or defaults")
		}
	}

	return &Config{
		Port:                   getEnv("GATEWAY_PORT", ":8080"),
		MessengerHost:          getEnv("MESSENGER_HOST", "http://localhost:8081"),
		TxHost:                 getEnv("TRANSACTION_HOST", "http://localhost:8083"),
		MediaHost:              getEnv("MEDIA_HOST", "http://localhost:8084"),
		JWTSecret:              readSecret("JWT_SECRET"),
		GatewaySignatureSecret: readSecret("GATEWAY_SIGNATURE_SECRET"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// readSecret returns the value of envName, or — if envName+"_FILE" is set —
// the contents of that file (trimmed). Lets dev use plain env vars and
// prod use Docker secrets without code changes.
func readSecret(envName string) string {
	if path := os.Getenv(envName + "_FILE"); path != "" {
		if data, err := os.ReadFile(path); err == nil {
			return strings.TrimSpace(string(data))
		}
	}
	return os.Getenv(envName)
}
