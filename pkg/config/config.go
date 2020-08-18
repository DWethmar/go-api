package config

import (
	"os"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBDriverName   string
	CreateDBScript string
}

func LoadConfig() Config {
	dbHost := os.Getenv("API_PQ_HOST")
	dbPort := os.Getenv("API_PG_PORT")
	dbUser := os.Getenv("API_PQ_USER")
	dbPassword := os.Getenv("API_PQ_PASSWORD")
	dbName := os.Getenv("API_PQ_DB_NAME")
	dbDriverName := os.Getenv("API_DRIVER_NAME")
	createDBScript := os.Getenv("CREATE_DB_SQL_FILE")

	return Config{
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		dbDriverName,
		createDBScript,
	}
}
