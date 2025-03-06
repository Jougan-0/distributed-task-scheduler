package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPass     string
	DBName     string
	DBPort     string
	ServerPort string
	SSLMode    string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "admin"),
		DBPass:     getEnv("DB_PASS", "admin"),
		DBName:     getEnv("DB_NAME", "tasks_db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
	}
	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort, c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
