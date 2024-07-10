package config

import (
	"os"
	"time"
)

type ENV string

const (
	ENVLocal ENV = "local"
	ENVProd  ENV = "prod"
)

func strToENV(env string) ENV {
	if e := ENV(env); e == ENVLocal || e == ENVProd {
		return e
	}
	return ENVProd
}

type Clickhouse struct {
	User     string
	Password string
	Host     string
	Port     string
}
type Server struct {
	Port    string
	Timeout time.Duration
}

type Config struct {
	ENV        ENV
	Clickhouse Clickhouse
	Server     Server
}

func getENV(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}

func Setup() Config {
	env := strToENV(getENV("ENV", "prod"))

	serverPort := ":" + getENV("SERVER_PORT", "8080")
	timeout, err := time.ParseDuration(getENV("TIMEOUT", "10s"))
	if err != nil {
		timeout = 10 * time.Second
	}

	clickhousePort := getENV("CLICKHOUSE_PORT", "9000")
	clickhouseHost := getENV("CLICKHOUSE_HOST", "localhost")
	clickhouseUser := getENV("CLICKHOUSE_ADMIN_USER", "clickhouse")
	clickhousePassword := getENV("CLICKHOUSE_ADMIN_PASSWORD", "clickhouse")

	return Config{
		ENV: env,
		Clickhouse: Clickhouse{
			User:     clickhouseUser,
			Password: clickhousePassword,
			Host:     clickhouseHost,
			Port:     clickhousePort,
		},
		Server: Server{
			Port:    serverPort,
			Timeout: timeout,
		},
	}
}
