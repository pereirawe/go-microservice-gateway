package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config is the structure that contains all the application configuration.
type Config struct {
	// APPConfigs
	APPName   string
	APPPort   string
	APPHost   string
	APPEnv    string
	APPUser   string
	APPPass   string
	APPSecret string

	// DBConfigs
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	APIKey     string

	// Microservices
	MSAi       string
	MSPayments string
	MSBi       string
	MSReports  string
}

// LoadConfig charge variables from .env file
// If a variable is not found, it will use the default value
func LoadConfig() (*Config, error) {
	// Load .env file
	// The `.` indicates that it will search for the .env file in the current directory.
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("No se encontró el archivo .env, usando valores por defecto.")
	}

	appName := setEnv("APP_NAME", "gateway")
	appPort := setEnv("APP_PORT", "8080")
	appHost := setEnv("APP_HOST", "localhost")
	appEnv := setEnv("APP_ENV", "development")
	appUser := setEnv("APP_USER", "admin")
	appPass := setEnv("APP_PASS", "password")
	appSecret := setEnv("APP_SECRET", "secret-key")

	dbHost := setEnv("DB_HOST", "localhost")
	dbPort := setEnv("DB_PORT", "5432")
	dbUser := setEnv("DB_USER", "postgres")
	dbPassword := setEnv("DB_PASSWORD", "password")
	dbName := setEnv("DB_NAME", "app_db")
	apiKey := setEnv("API_KEY", "default-api-key")

	// Microservices
	msAiIp := setEnv("MS_AI_IP", "localhost")
	msAiPort := setEnv("MS_AI_PORT", "8081")
	msPaymentsIp := setEnv("MS_PAYMENTS_IP", "localhost")
	msPaymentsPort := setEnv("MS_PAYMENTS_PORT", "8082")
	msBiIp := setEnv("MS_BI_IP", "localhost")
	msBiPort := setEnv("MS_BI_PORT", "8083")
	msReportsIp := setEnv("MS_REPORTS_IP", "localhost")
	msReportsPort := setEnv("MS_REPORTS_PORT", "8084")

	// Return the configuration in the Config structure
	cfg := &Config{

		APPName:   appName,
		APPPort:   appPort,
		APPHost:   appHost,
		APPEnv:    appEnv,
		APPUser:   appUser,
		APPPass:   appPass,
		APPSecret: appSecret,

		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		APIKey:     apiKey,

		MSAi:       msAiIp + ":" + msAiPort,
		MSPayments: msPaymentsIp + ":" + msPaymentsPort,
		MSBi:       msBiIp + ":" + msBiPort,
		MSReports:  msReportsIp + ":" + msReportsPort,
	}

	return cfg, nil
}

// SetEnv sets the value of an environment variable if it is not already set.
// If the variable is already set, it returns the current value.
// If the variable is not set, it sets the value to the default value.
func setEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
