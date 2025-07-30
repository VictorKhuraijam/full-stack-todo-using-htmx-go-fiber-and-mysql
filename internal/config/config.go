package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct{
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	Port string
}

func Load() *Config  {
	godotenv.Load()

	return &Config{
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "3306"),
		DBUser: getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName: getEnv("DB_NAME", "todo_db"),
		Port: getEnv("PORT", "3000"),
	}
}

func getEnv(key, defaultValue string) string  {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
