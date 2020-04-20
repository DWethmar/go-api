package config

import (
	"os"
)

type AuthEnv struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBDriverName string
}

func LoadAuthEnv() AuthEnv {
	dbHost := os.Getenv("AUTH_PQ_HOST")
	dbPort := os.Getenv("AUTH_PG_PORT")
	dbUser := os.Getenv("AUTH_PQ_USER")
	dbPassword := os.Getenv("AUTH_PQ_PASSWORD")
	dbName := os.Getenv("AUTH_PQ_DB_NAME")
	dbDriverName := os.Getenv("AUTH_DRIVER_NAME")

	return AuthEnv{
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		dbDriverName,
	}
}
