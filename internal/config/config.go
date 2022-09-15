package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	// SwdApp the app title
	SwdApp = "swd-app"
)

// New config
func New(app string) *Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err.Error())
	}
	cfg := new(Config)

	httpServerHost := os.Getenv("HTTP_SERVER_PORT")

	if httpServerHost == "" {
		httpServerHost = "8060"
	}

	switch app {
	case SwdApp:
		cfg.HTTPServerConfig = &HTTPServerConfig{
			HTTPServerHost: httpServerHost,
			WaitToShutdown: 5,
		}
		cfg.MySQLRepositoryConfig = &MySQLRepositoryConfig{
			DBUser:     os.Getenv("MYSQL_USER"),
			DBPassword: os.Getenv("MYSQL_PASSWORD"),
			DBName:     os.Getenv("MYSQL_DB_NAME"),
			DBHost:     os.Getenv("MYSQL_HOST"),
			DBPort:     os.Getenv("MYSQL_DB_PORT"),
		}
	}

	return cfg
}
