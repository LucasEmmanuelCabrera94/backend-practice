package config

import (
	"fmt"
	"os"
)

type Config struct {
	MySQLDSN    string
	AppPort     string
	JWTSecret   string
}

func Load() *Config {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		user := getEnv("MYSQL_USER", "backend")
		pass := getEnv("MYSQL_PASSWORD", "backendpw")
		dbName := getEnv("MYSQL_DATABASE", "backend")
		host := getEnv("MYSQL_HOST", "")
		port := getEnv("MYSQL_PORT", "")

		if host == "" {
			if os.Getenv("DOCKER") == "1" {
				host = "mysql"
			} else {
				host = "127.0.0.1"
			}
		}

		if port == "" {
			if host == "127.0.0.1" {
				port = "3307"
			} else {
				port = "3306"
			}
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbName)
	}

	return &Config{
		MySQLDSN:  dsn,
		AppPort:   getEnv("APP_PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET_KEY", "super_secret_key"),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
