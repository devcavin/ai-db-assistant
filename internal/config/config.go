package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	minJWTSecretLength = 32
	placeholderPrefix = "CHANGE_ME"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiry      time.Duration
	Environment    string
	ProxyServerURL string
}

// Every problem found is collected and returned in a single error so that all of them can be fixed in one go.
func Load() (*Config, error) {
	var problems []string

	required := func(key string) string {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			problems = append(problems, key+" is required")
		}
		return value
	}

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    required("DATABASE_URL"),
		JWTSecret:      required("JWT_SECRET"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		ProxyServerURL: getEnv("PROXY_SERVER_URL", "http://localhost:8000"),
	}

	if cfg.JWTSecret != "" {
		if strings.HasPrefix(cfg.JWTSecret, placeholderPrefix) {
			problems = append(problems, "JWT_SECRET still has the placeholder value from .env.example")
		} else if len(cfg.JWTSecret) < minJWTSecretLength {
			problems = append(problems, fmt.Sprintf("JWT_SECRET must be at least %d characters", minJWTSecretLength))
		}
	}

	expiry, err := getDurationEnv("JWT_EXPIRY", 24*time.Hour)
	if err != nil {
		problems = append(problems, err.Error())
	}
	cfg.JWTExpiry = expiry

	if len(problems) > 0 {
		return nil, fmt.Errorf("invalid configuration: %s", strings.Join(problems, "; "))
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue, fmt.Errorf("%s is not a valid duration (e.g. 24h, 90m)", key)
	}
	if duration <= 0 {
		return defaultValue, fmt.Errorf("%s must be greater than zero", key)
	}

	return duration, nil
}