package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	host       string
	port       string
	user       string
	password   string
	dbName     string
	driverName string
}

func LoadEnv() Config {
	host := os.Getenv("PQ_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PQ_USER")
	password := os.Getenv("PQ_PASSWORD")
	dbName := os.Getenv("PQ_DB_NAME")
	driverName := os.Getenv("DRIVER_NAME")

	return Config{
		host,
		port,
		user,
		password,
		dbName,
		driverName,
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
	return config.driverName, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.host,
		config.port,
		config.user,
		config.password,
		config.dbName,
	)
}
