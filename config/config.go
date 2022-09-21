package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type dbConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}

type AppConfig struct {
	BaseURL            string
	Env                string
	AppName            string
	JWTSecret          []byte
	JWTExpireInMinutes int
	StoragePath        string
	DBConfig           dbConfig
}

func getEnv(key, defaultVal string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	env := os.Getenv(key)
	if env == "" {
		return defaultVal
	}
	return env
}

var Config = AppConfig{
	BaseURL:            getEnv("BASE_URL", "http://localhost:8080"),
	Env:                getEnv("ENV", "testing"),
	AppName:            "kamari",
	JWTSecret:          []byte("kamari-secret"),
	JWTExpireInMinutes: 15,
	DBConfig: dbConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "kamari_db"),
		Port:     getEnv("DB_PORT", "5432"),
		SSLMode:  getEnv("DB_SSL", "disable"),
	},
}
