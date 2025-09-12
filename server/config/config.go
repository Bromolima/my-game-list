package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var Env Environment

func LoadEnvironment() error {
	err := godotenv.Load()
	if err != nil {
		slog.Error("failed to load environment variables", slog.String("error", err.Error()))
		return err
	}

	Env = Environment{
		Env:       getEnv("ENV", "development"),
		ApiPort:   getEnv("API_PORT", "8080"),
		SecretKey: getEnv("JWT_SECRET_KEY", ""),
		DB: Mysql{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "root"),
			Name:     getEnv("DB_NAME", "url-shortener"),
		},
	}

	slog.Info("environment variables loaded successfully")

	return nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaultValue
}
