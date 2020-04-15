package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBDriverName   string
	CreateDBScript string
}

func LoadEnv() Env {
	dbHost := os.Getenv("API_PQ_HOST")
	dbPort := os.Getenv("API_PG_PORT")
	dbUser := os.Getenv("API_PQ_USER")
	dbPassword := os.Getenv("API_PQ_PASSWORD")
	dbName := os.Getenv("API_PQ_DB_NAME")
	dbDriverName := os.Getenv("API_DRIVER_NAME")
	createDBScript := os.Getenv("CREATE_DB_SQL_FILE")

	return Env{
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

func LoadEnvFile(path ...string) Env {
	if path == nil {
		path = []string{defaultEnvPath}
	}
	err := godotenv.Load(path...)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return LoadEnv()
}

func GetPostgresConnectionInfo(config Env) (string, string) {
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
