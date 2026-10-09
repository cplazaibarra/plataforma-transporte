package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port         string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	UploadDir    string
	AllowedOrigs string
	JWTSecret    string
}

func LoadConfig() *Config {
	port := getEnv("PORT", "8095")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5438")
	dbUser := getEnv("DB_USER", "lastmile_user")
	dbPass := getEnv("DB_PASSWORD", "lastmile_password")
	dbName := getEnv("DB_NAME", "lastmile_db")
	dbSSL := getEnv("DB_SSLMODE", "disable")
	uploadDir := getEnv("UPLOAD_DIR", "./uploads")
	allowedOrigs := getEnv("ALLOWED_ORIGINS", "*")
	jwtSecret := getEnv("JWT_SECRET", "lastmile_dev_jwt_secret_key_32bytes!!")

	return &Config{
		Port:         port,
		DBHost:       dbHost,
		DBPort:       dbPort,
		DBUser:       dbUser,
		DBPassword:   dbPass,
		DBName:       dbName,
		DBSSLMode:    dbSSL,
		UploadDir:    uploadDir,
		AllowedOrigs: allowedOrigs,
		JWTSecret:    jwtSecret,
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
