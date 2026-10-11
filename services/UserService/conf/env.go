package conf

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	GRPCPort    string
	Environment string
	HTTPPort    string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, relying on system env vars")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "messenger_users"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		//		GRPCPort:    getEnv("GRPC_PORT", "50051"),
		Environment: getEnv("ENVIRONMENT", "dev"),
		HTTPPort:    getEnv("HTTPPORT", "50051"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
