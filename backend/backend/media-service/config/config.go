package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	MediaPort        string
	MediaGRPCPort    string
	MediaDatabaseURL string

	S3Endpoint       string
	S3PublicEndpoint string
	S3AccessKey      string
	S3SecretKey      string
	S3UseSSL         bool
	S3Bucket         string

	MessengerGRPCAddr      string
	GatewaySignatureSecret string

	// Presigned URL lifetimes
	PresignUploadExpiry   time.Duration
	PresignViewExpiry     time.Duration
	PresignDownloadExpiry time.Duration
}

func Load() *Config {
	return &Config{
		MediaPort:        getEnv("MEDIA_PORT", ":8084"),
		MediaGRPCPort:    getEnv("MEDIA_GRPC_PORT", ":9094"),
		MediaDatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/suda?sslmode=disable"),

		S3Endpoint:       getEnv("S3_ENDPOINT", "localhost:9000"),
		S3PublicEndpoint: getEnv("S3_PUBLIC_ENDPOINT", "localhost:9000"),
		S3AccessKey:      getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey:      getEnv("S3_SECRET_KEY", "minioadmin"),
		S3UseSSL:         getBoolEnv("S3_USE_SSL", false),
		S3Bucket:         getEnv("S3_DEFAULT_BUCKET", "media"),

		MessengerGRPCAddr:      getEnv("MESSENGER_GRPC_ADDR", "localhost:9091"),
		GatewaySignatureSecret: readSecret("GATEWAY_SIGNATURE_SECRET"),

		PresignUploadExpiry:   getDurationEnv("PRESIGN_UPLOAD_EXPIRY", 15*time.Minute),
		PresignViewExpiry:     getDurationEnv("PRESIGN_VIEW_EXPIRY", 4*time.Hour),
		PresignDownloadExpiry: getDurationEnv("PRESIGN_DOWNLOAD_EXPIRY", 1*time.Hour),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		return val == "true" || val == "1"
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if val, ok := os.LookupEnv(key); ok {
		if seconds, err := strconv.Atoi(val); err == nil {
			return time.Duration(seconds) * time.Second
		}
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
