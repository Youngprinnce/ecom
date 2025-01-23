package config

import (
	"os"

	"github.com/joho/godotenv"
)

var Envs = initConfig()

type Config struct {
	DB DB
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
	}
}

func getEnvOrPanic(key, err string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(err)
}
