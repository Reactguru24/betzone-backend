package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port            string
	Environment     string
	BetkraftBaseURL string
	BetkraftAPIKey  string
	BetkraftAppKey  string
	JWTSecret       string
	CallbackURL     string // e.g., "http://localhost:8080/api/v1/callbacks"
	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() (*Config, error) {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		BetkraftBaseURL: getEnv("BETKRAFT_BASE_URL", "https://api.staging.betkraft.co.uk"),
		BetkraftAPIKey:  getEnv("API_KEY", ""),
		BetkraftAppKey:  getEnv("APP_KEY", ""),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		CallbackURL:     getEnv("CALLBACK_URL", "http://localhost:8080/api/v1/callbacks"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "3306"),
		DBUser:          getEnv("DB_USER", "betzone"),
		DBPassword:      getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "betzone"),
	}, nil
}

// GetDatabaseDSN returns the MySQL DSN string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

// GetCallbackURLs returns all callback URLs as a map
func (c *Config) GetCallbackURLs() map[string]string {
	return map[string]string{
		"player_info": c.CallbackURL + "/player_info",
		"bet":         c.CallbackURL + "/bet",
		"win":         c.CallbackURL + "/win",
		"rollback":    c.CallbackURL + "/rollback",
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
