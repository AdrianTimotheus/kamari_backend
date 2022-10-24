package config

import (
	_ "github.com/joho/godotenv/autoload"
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

type redisConfig struct {
	Host     string
	Password string
}

type seaLabsPayConfig struct {
	MerchantName string
	MerchantCode string
	ApiKey       string
}

type courierService struct {
	Url string
	Key string
}

type AppConfig struct {
	BaseURL            string
	Env                string
	AppName            string
	JWTSecret          []byte
	JWTExpireInMinutes int
	StoragePath        string
	DBConfig           dbConfig
	RedisConfig        redisConfig
	SeaLabsPayConfig   seaLabsPayConfig
	CourierService     courierService
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
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "kamari_db"),
		Port:     getEnv("DB_PORT", "5432"),
		SSLMode:  getEnv("DB_SSL", "disable"),
	},
	RedisConfig: redisConfig{
		Host:     getEnv("REDIS_HOST", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
	},
	SeaLabsPayConfig: seaLabsPayConfig{
		MerchantName: getEnv("MERCHANT_NAME", ""),
		MerchantCode: getEnv("MERCHANT_CODE", ""),
		ApiKey:       getEnv("SEA_LABS_PAY_API_KEY", ""),
	},
	CourierService: courierService{
		Url: getEnv("COURIER_SERVICE_URL", ""),
		Key: getEnv("COURIER_SERVICE_KEY", ""),
	},
}
