package config

import (
	"os"
)

type Config struct {
	Port              string
	Environment       string
	BetkraftBaseURL   string
	BetkraftAPIKey    string
	BetkraftAppKey    string
}

func LoadConfig() (*Config, error) {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		BetkraftBaseURL: getEnv("BETKRAFT_BASE_URL", "https://api.staging.betkraft.co.uk"),
		BetkraftAPIKey:  getEnv("API_KEY", ""),
		BetkraftAppKey:  getEnv("APP_KEY", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
