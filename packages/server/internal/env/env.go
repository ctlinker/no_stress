package env

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	PORT               string
	JWT_REFRESH_SECRET string
	JWT_ACCESS_SECRET  string
	DATABASE_URL       string
}

func Load() Config {
	return Config{
		PORT:               GetOrWarn("SERVER_PORT"),
		JWT_ACCESS_SECRET:  GetOrWarn("JWT_ACCESS_SECRET"),
		JWT_REFRESH_SECRET: GetOrWarn("JWT_REFRESH_SECRET"),
		DATABASE_URL:       GetOrWarn("DATABASE_URL"),
	}
}

func GetOrWarn(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Println(fmt.Errorf("!!! Warning environement var `%s` not found", name))
	}
	return v
}

func GetOrDefault(name string, def string) string {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	return v
}
