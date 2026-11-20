package config

import (
	"os"
	"time"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiry     time.Duration
	Environment   string
	ProxyServerURL string
}


func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiry:      getDurationEnv("JWT_EXPIRY", time.Hour * 24), // 24 hours
		Environment:    getEnv("ENVIRONMENT", "development"),
		ProxyServerURL: getEnv("PROXY_SERVER_URL", "http://localhost:8000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
    value, exists := os.LookupEnv(key)
    if !exists {
        return defaultValue
    }

    duration, err := time.ParseDuration(value)
    if err != nil {
        return defaultValue
    }

    return duration
}