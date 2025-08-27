package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	// Server Config
	Port string

	// Database Config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis Config
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string

	// JWT Config
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExpiry  string // in hours
	JWTRefreshExpiry string // in hours

}

func GetConfig() *Config {
	cfg := &Config{
		// Server
		Port: getEnv("PORT"),

		// Database
		DBHost:     getEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME"),

		// Redis
		RedisHost:     getEnv("REDIS_HOST"),
		RedisPort:     getEnv("REDIS_PORT"),
		RedisPassword: getEnv("REDIS_PASSWORD"),
		RedisDB:       getEnv("REDIS_DB"),

		// JWT
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET"),
		JWTAccessExpiry:  getEnv("JWT_ACCESS_EXPIRY"),
		JWTRefreshExpiry: getEnv("JWT_REFRESH_EXPIRY"),
	}

	// Validate all required environment variables
	ValidateConfig(cfg)

	return cfg
}

// GetDatabaseDSN returns MySQL DSN string with UTC timezone
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func ValidateConfig(cfg *Config) {
	// Required environment variables (DB_PASSWORD is optional - can be blank)
	requiredEnvVars := map[string]string{
		"PORT":               cfg.Port,
		"DB_HOST":            cfg.DBHost,
		"DB_PORT":            cfg.DBPort,
		"DB_USER":            cfg.DBUser,
		"DB_NAME":            cfg.DBName,
		"REDIS_HOST":         cfg.RedisHost,
		"REDIS_PORT":         cfg.RedisPort,
		"REDIS_DB":           cfg.RedisDB,
		"JWT_ACCESS_SECRET":  cfg.JWTAccessSecret,
		"JWT_REFRESH_SECRET": cfg.JWTRefreshSecret,
		"JWT_ACCESS_EXPIRY":  cfg.JWTAccessExpiry,
		"JWT_REFRESH_EXPIRY": cfg.JWTRefreshExpiry,
	}

	var missingEnvVars []string
	for envVar, value := range requiredEnvVars {
		if value == "" {
			missingEnvVars = append(missingEnvVars, envVar)
		}
	}

	if len(missingEnvVars) > 0 {
		log.Printf("Missing required environment variables: %v", missingEnvVars)
		log.Fatal("Please set all required environment variables in .env file")
	}

	// Validate JWT expiry values are valid integers
	accessExpiry, err := strconv.Atoi(cfg.JWTAccessExpiry)
	if err != nil {
		log.Printf("Invalid JWT_ACCESS_EXPIRY value '%s': must be a valid integer (hours)", cfg.JWTAccessExpiry)
		log.Fatal("JWT_ACCESS_EXPIRY must be a valid integer representing hours")
	}
	if accessExpiry <= 0 {
		log.Printf("Invalid JWT_ACCESS_EXPIRY value '%d': must be greater than 0", accessExpiry)
		log.Fatal("JWT_ACCESS_EXPIRY must be a positive integer")
	}

	refreshExpiry, err := strconv.Atoi(cfg.JWTRefreshExpiry)
	if err != nil {
		log.Printf("Invalid JWT_REFRESH_EXPIRY value '%s': must be a valid integer (hours)", cfg.JWTRefreshExpiry)
		log.Fatal("JWT_REFRESH_EXPIRY must be a valid integer representing hours")
	}
	if refreshExpiry <= 0 {
		log.Printf("Invalid JWT_REFRESH_EXPIRY value '%d': must be greater than 0", refreshExpiry)
		log.Fatal("JWT_REFRESH_EXPIRY must be a positive integer")
	}
}
