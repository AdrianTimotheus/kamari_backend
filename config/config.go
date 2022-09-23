package config

import (
<<<<<<< HEAD
<<<<<<< HEAD
=======
	_ "github.com/joho/godotenv/autoload"
>>>>>>> feat:Get promotion img URL for <img> src sorted by latest promotion
=======
	_ "github.com/joho/godotenv/autoload"
>>>>>>> 389454c5afc0ee87bf80a2c41600fe616111f407
	"os"

	_ "github.com/joho/godotenv/autoload"
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
