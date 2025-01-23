package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Envs = initConfig()

type Config struct {
	DB DB
	PORT string
	JWT_EXPIRE_IN_SECONDS int64
	JWT_SECRET string
}

type DB struct {
	User string
	Passwd string
	Net string
	Addr string
	DBName string
	AllowNativePasswords bool
	ParseTime bool
}

func initConfig() Config {
	godotenv.Load()
	return Config{
		DB: DB{
			User: getEnvOrPanic("DB_USER", "DB_USER is required"),
			Passwd: getEnvOrPanic("DB_PASSWD", "DB_PASSWD is required"),
			Net: getEnvOrPanic("DB_NET", "DB_NET is required"),
			Addr: getEnvOrPanic("DB_ADDR", "DB_ADDR is required"),
			DBName: getEnvOrPanic("DB_NAME", "DB_NAME is required"),
			AllowNativePasswords: true,
			ParseTime: true,
		},
		PORT: getEnvOrPanic("PORT", "PORT is required"),
		JWT_EXPIRE_IN_SECONDS: getEnvAsInt("JWT_EXPIRE_IN_SECONDS", 3600 * 24 * 7),
		JWT_SECRET: getEnvOrPanic("JWT_SECRET", "JWT_SECRET is required"),
	}
}

func getEnvOrPanic(key, err string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(err)
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}
	return fallback
}
