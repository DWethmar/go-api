package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
)

var defaultEnvPath = ".env"

func LoadEnvFile(path ...string) Config {
	if path == nil {
		path = []string{defaultEnvPath}
	}
	err := godotenv.Load(path...)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return LoadConfig()
}

func GetPostgresConnectionInfo(env Config) (string, string) {
	cParts := []string{
		fmt.Sprintf("host=%s", env.DBHost),
		fmt.Sprintf("port=%s", env.DBPort),
		"sslmode=disable",
	}

	if env.DBUser != "" {
		cParts = append(cParts, fmt.Sprintf("user=%s", env.DBUser))
	}

	if env.DBPassword != "" {
		cParts = append(cParts, fmt.Sprintf("password=%s", env.DBPassword))
	}

	if env.DBName != "" {
		cParts = append(cParts, fmt.Sprintf("dbname=%s", env.DBName))
	}

	return env.DBDriverName, strings.Join(cParts, " ")
}
