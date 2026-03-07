package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	AppEnv     string
}

// LoadConfig reads the .env file and maps environment variables into a Config struct.
func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		// .env is optional; continue with existing env
		_ = err
	}

	return Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "hsmart_saas"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		AppEnv:     getEnv("APP_ENV", "development"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
