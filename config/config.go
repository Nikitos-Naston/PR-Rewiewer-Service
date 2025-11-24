package config

import (
	"PRreviewService/internal/messages"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	configPath string = "config/.env"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
}

func LoadConfig() Config {

	err := godotenv.Load(configPath)

	if err != nil {
		messages.SendLogMessage(logrus.ErrorLevel, "File .env was not found ", err)
	}
	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "PRreviewDB"),
		Port:       getEnv("PORT", "8080"),
	}
}

func getEnv(val string, defaultVal string) string {
	if value, exists := os.LookupEnv(val); exists {
		return value
	}
	return defaultVal
}
