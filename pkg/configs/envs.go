package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	ProductionEnv      = "production"
	DatabaseTimeout    = 5 * time.Second
	ProductCachingTime = 1 * time.Minute
)

type Config struct {
	URI           string
	Port          string
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	AuthSecret    string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		URI: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			getEnv("DB_USERNAME", "postgres"),
			getEnv("DB_PASSWORD", "mypassword"),
			getEnv("DB_HOST", "127.0.0.1"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_DATABASE", "db_washit"),
		),
		Port:          getEnv("APP_PORT", "8080"),
		DBUser:        getEnv("DB_USERNAME", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "mypassword"),
		DBHost:        getEnv("DB_HOST", "127.0.0.1"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBName:        getEnv("DB_DATABASE", "db_washit"),
		RedisHost:     getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
		AuthSecret:    getEnv("AUTH_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return int(i)
	}

	return fallback
}
