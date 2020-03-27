package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBDriverName string
}

func LoadEnv() Config {
	dbHost := os.Getenv("PQ_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbUser := os.Getenv("PQ_USER")
	dbPassword := os.Getenv("PQ_PASSWORD")
	dbName := os.Getenv("PQ_DB_NAME")
	dbDriverName := os.Getenv("DRIVER_NAME")

	return Config{
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		dbDriverName,
	}
}

var defaultEnvPath = "../.env"

func LoadEnvFile(path ...string) Config {
	if path == nil {
		path = []string{defaultEnvPath}
	}
	err := godotenv.Load(path...)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return LoadEnv()
}

func GetPostgresConnectionInfo(config Config) (string, string) {
	return config.DBDriverName, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)
}
