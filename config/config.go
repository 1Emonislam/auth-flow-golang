package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBPort         string
	BitcoinRPCURL  string
	BitcoinRPCUser string
	BitcoinRPCPass string
	ServerPort     string
	WebhookSecret  string
	JWTSecret      string
	BaseURL        string
	BitcoinNetwork string
	// Redis configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	TokenExpiry   int // Token expiry in hours
	// Email configuration
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
	SMTPFromName string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		DBPort:         os.Getenv("DB_PORT"),
		BitcoinRPCURL:  os.Getenv("BITCOIN_RPC_URL"),
		BitcoinRPCUser: os.Getenv("BITCOIN_RPC_USER"),
		BitcoinRPCPass: os.Getenv("BITCOIN_RPC_PASS"),
		BitcoinNetwork: os.Getenv("BITCOIN_NETWORK"), // e.g., "mainnet", "testnet", "regtest"
		ServerPort:     os.Getenv("SERVER_PORT"),
		WebhookSecret:  os.Getenv("WEBHOOK_SECRET"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		BaseURL:        os.Getenv("BASE_URL"),
		// Redis configuration
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       0, // Default DB
		TokenExpiry:   24, // Default 24 hours
		// Email configuration
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),
		SMTPFromName: os.Getenv("SMTP_FROM_NAME"),
	}
}
