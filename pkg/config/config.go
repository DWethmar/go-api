package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
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

func LoadEnv() Config {
	dbHost := os.Getenv("PQ_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbUser := os.Getenv("PQ_USER")
	dbPassword := os.Getenv("PQ_PASSWORD")
	dbName := os.Getenv("PQ_DB_NAME")
	dbDriverName := os.Getenv("DRIVER_NAME")
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
	cParts := []string{
		fmt.Sprintf("host=%s", config.DBHost),
		fmt.Sprintf("port=%s", config.DBPort),
		"sslmode=disable",
	}

	if config.DBUser != "" {
		cParts = append(cParts, fmt.Sprintf("user=%s", config.DBUser))
	}

	if config.DBPassword != "" {
		cParts = append(cParts, fmt.Sprintf("password=%s", config.DBPassword))
	}

	if config.DBName != "" {
		cParts = append(cParts, fmt.Sprintf("dbname=%s", config.DBName))
	}

	return config.DBDriverName, strings.Join(cParts, " ")
}
