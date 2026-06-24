package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                string
	GRPCPort               string
	DatabaseURL            string
	RedisAddr              string
	RedisPass              string
	AppSecret              string
	ContextTimeout         time.Duration
	TransactionServiceAddr string
	MediaServiceAddr       string
	InitDataSignToken      string //TODO ИЗМЕНИТЬ
	GatewaySignatureSecret string
	// MessageEncryptionKey — 64-char hex (AES-256) for at-rest message content
	// encryption. Empty disables encryption (messages stored as plaintext).
	MessageEncryptionKey string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println(" No .env file found, using system envs")
		}
	}

	
	return &Config{
		AppPort:                getEnv("MESSENGER_PORT", ":8081"),
		GRPCPort:               getEnv("MESSENGER_GRPC_PORT", ":9081"),
		DatabaseURL:            getEnv("DATABASE_URL", "postgres://postgres:123@localhost:5432/suda?sslmode=disable"),
		RedisAddr:              getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:              getEnv("REDIS_PASSWORD", ""),
		AppSecret:              getEnv("APP_SECRET", ""),
		ContextTimeout:         5 * time.Second, //TODO Хардкод таймаута, можно оставить не критично но не забывать про него
		TransactionServiceAddr: getEnv("TRANSACTION_SERVICE_ADDR", "localhost:9082"),
		MediaServiceAddr:       getEnv("MEDIA_SERVICE_ADDR", "localhost:9094"),
		InitDataSignToken:      readSecret("INIT_DATA_SIGN_TOKEN"),
		GatewaySignatureSecret: readSecret("GATEWAY_SIGNATURE_SECRET"),
		MessageEncryptionKey:   readSecret("MESSAGE_ENCRYPTION_KEY"),
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
