package config

import (
	"os"
)

// Config collection of the config variables.
type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBDriverName       string
	MigrationFiles     string
	DBMigrationVersion int
	TestWithDB         bool
}

// Load collects the necessary env vars and returns them in a struct.
func Load() Config {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")

	dbName := "go-api"
	if v, success := os.LookupEnv("POSTGRES_DATABASE"); success {
		dbName = v
	}

	dbDriverName := os.Getenv("DB_DRIVER_NAME")

	migrationFiles := "/app/migrations"
	if v, success := os.LookupEnv("MIGRATION_FILES"); success {
		migrationFiles = v
	}

	testWithDB := false
	if v, success := os.LookupEnv("TEST_WITH_DB"); success {
		testWithDB = success && v == "true"
	}

	dbMigrationVersion := 1

	return Config{
		DBHost:             dbHost,
		DBPort:             dbPort,
		DBUser:             dbUser,
		DBPassword:         dbPassword,
		DBName:             dbName,
		DBDriverName:       dbDriverName,
		MigrationFiles:     migrationFiles,
		DBMigrationVersion: dbMigrationVersion,
		TestWithDB:         testWithDB,
	}
}
