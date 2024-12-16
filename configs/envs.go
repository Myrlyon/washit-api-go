package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DSN        string
	AuthSecret string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "mypassword"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "db_washit"),
		AuthSecret: getEnv("AUTH_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
