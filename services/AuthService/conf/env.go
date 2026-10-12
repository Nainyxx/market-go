package conf

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	GRPCPort       string
	Environment    string
	HTTPPort       string
	JWTSecret      string
	AccessTTL      time.Duration
	UserServiceURL string
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
		Environment:    getEnv("ENVIRONMENT", "dev"),
		HTTPPort:       getEnv("HTTPPORT", "50051"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		AccessTTL:      getDuration("ACCESS_TTL", 15*time.Minute),
		UserServiceURL: getEnv("USER_SERVICE_URL", "http://localhost:8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("invalid %s=%q, using default %s", key, value, defaultValue)
		return defaultValue
	}
	return d
}
