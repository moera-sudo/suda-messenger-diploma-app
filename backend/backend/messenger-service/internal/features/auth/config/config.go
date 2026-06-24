package config

import (
	"os"
	"strings"
	"time"
)

type AuthConfig struct {
	JWTSecret string
	AccessTokenTTL time.Duration
	RefreshTokenTTL time.Duration

}

func Load() *AuthConfig {
	return &AuthConfig{
		JWTSecret: readSecret("JWT_SECRET"),

		AccessTokenTTL: getEnvDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL: getEnvDuration("REFRESH_TOKEN_TTL", 30*24*time.Hour),
	}
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val, exists := os.LookupEnv(key); exists {
		// time.ParseDuration understands "15m", "1h", "10s", ...
		if d, err := time.ParseDuration(val); err == nil {
			return d
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
