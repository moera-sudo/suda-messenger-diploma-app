package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	GRPCPort       string
	DatabaseURL    string
	RedisAddr      string
	RedisPass      string
	ContextTimeout time.Duration

	WalletEncryptionKey string // Симметричный ключ(hex, 32 байта) для AES-256-GCM шифрования приватных ключей кошельков юзеров

	BesuRPCURL         string
	ChainID            int64
	DeployerPrivateKey string

	SudaTokenAddress       string
	SudaNFTAddress         string
	SudaEscrowAddress      string
	SudaFundraisingAddress string
	SudaMarketplaceAddress string

	MessengerGRPCAddr      string
	GatewaySignatureSecret string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println(" No .env file found, using system envs")
		}
	}

	return &Config{
		AppPort:        getEnv("TRANSACTION_PORT", ":8082"),
		GRPCPort:       getEnv("TRANSACTION_GRPC_PORT", ":9082"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:123@localhost:5432/suda?sslmode=disable"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:      getEnv("REDIS_PASSWORD", ""),
		ContextTimeout: 5 * time.Second, //TODO хардкод, не забыть вынести

		WalletEncryptionKey: readSecret("WALLET_ENCRYPTION_KEY"),

		BesuRPCURL:         getEnv("BESU_RPC_URL", "http://localhost:8545"),
		ChainID:            getEnvInt64("BESU_CHAIN_ID", 1337),
		DeployerPrivateKey: readSecret("ETH_PRIVATE_KEY"),

		SudaTokenAddress:       getEnv("SUDA_TOKEN_ADDRESS", ""),
		SudaNFTAddress:         getEnv("SUDA_NFT_ADDRESS", ""),
		SudaEscrowAddress:      getEnv("SUDA_ESCROW_ADDRESS", ""),
		SudaFundraisingAddress: getEnv("SUDA_FUNDRAISING_ADDRESS", ""),
		SudaMarketplaceAddress: getEnv("SUDA_MARKETPLACE_ADDRESS", ""),

		MessengerGRPCAddr:      getEnv("MESSENGER_GRPC_ADDR", "localhost:9081"),
		GatewaySignatureSecret: readSecret("GATEWAY_SIGNATURE_SECRET"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
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
