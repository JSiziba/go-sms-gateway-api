package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBHost                  string
	DBPort                  int
	DBUser                  string
	DBPassword              string
	DBName                  string
	SSLMode                 string
	SSLRootCert             string
	ServerPort              int
	XRequireWhiskAuth       bool
	XRequireWhiskAuthSecret string
}

func LoadConfig() (config Config, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	config.DBHost = "localhost"
	config.DBPort = 5432
	config.DBUser = "postgres"
	config.SSLMode = "disable"
	config.SSLRootCert = ""
	config.DBPassword = ""
	config.DBName = "default"
	config.ServerPort = 8080
	config.XRequireWhiskAuth = false
	config.XRequireWhiskAuthSecret = ""

	if os.Getenv("DB_HOST") != "" {
		config.DBHost = os.Getenv("DB_HOST")
	}

	if os.Getenv("DB_PORT") != "" {
		config.DBPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	}

	if os.Getenv("DB_USER") != "" {
		config.DBUser = os.Getenv("DB_USER")
	}

	if os.Getenv("DB_PASSWORD") != "" {
		config.DBPassword = os.Getenv("DB_PASSWORD")
	}

	if os.Getenv("DB_NAME") != "" {
		config.DBName = os.Getenv("DB_NAME")
	}

	if os.Getenv("DB_SSLMODE") != "" {
		config.SSLMode = os.Getenv("DB_SSLMODE")
	}

	if os.Getenv("DB_SSLROOTCERT") != "" {
		config.SSLRootCert = os.Getenv("DB_SSLROOTCERT")
	}

	if os.Getenv("SERVER_PORT") != "" {
		config.ServerPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	}

	if os.Getenv("X_REQUIRE_WHISK_AUTH") != "" {
		config.XRequireWhiskAuth = strings.ToLower(os.Getenv("X_REQUIRE_WHISK_AUTH")) == "true"
	}
	if os.Getenv("X_REQUIRE_WHISK_AUTH_SECRET") != "" {
		config.XRequireWhiskAuthSecret = os.Getenv("X_REQUIRE_WHISK_AUTH_SECRET")
	}
	return config, nil
}

func (c Config) GetDBConnString() string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.SSLMode)

	if c.SSLRootCert != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", c.SSLRootCert)
	}
	return dsn
}
