package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/dwethmar/go-api/pkg/config"
)

func CreateDatabase(db *sql.DB, name string) {
	fmt.Printf("Creating database: %v \n", name)
	_, err := db.Exec("CREATE DATABASE " + name + ";")

	if err != nil {
		panic(err)
	}
}

func DropDatabase(db *sql.DB, name string) {
	fmt.Printf("Dropping database: %v \n", name)
	_, err := db.Exec("DROP DATABASE " + name)

	if err != nil {
		panic(err)
	}
}

func ExecSQLFile(db *sql.DB, sqlFile string) {

	fmt.Printf("Reading SQL file: %v \n", sqlFile)

	b, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatal(err)
	}

	sql := string(b)

	fmt.Printf("Excecuting SQL file: %v \n", sqlFile)

	_, err = db.Exec(sql)
	if err != nil {
		panic(err)
	}
}

// NewTestDB create new testy db. Returns a cleanup function and error.
func NewTestDB(c config.Config) (*sql.DB, error) {
	var db *sql.DB

	if c.DBDriverName == "" {
		return nil, errors.New("Config is missing database connection information: DBDriverName")
	}

	if c.DBHost == "" {
		return nil, errors.New("Config is missing database connection information: DBHost")
	}

	if c.DBName == "" {
		return nil, errors.New("Config is missing database connection information: DBName")
	}

	if c.DBPassword == "" {
		return nil, errors.New("Config is missing database connection information: DBPassword")
	}

	if c.DBPort == "" {
		return nil, errors.New("Config is missing database connection information: DBPort")
	}

	if c.DBUser == "" {
		return nil, errors.New("Config is missing database connection information: DBUser")
	}

	if c.MigrationFiles == "" {
		return nil, errors.New("Config is missing database connection information: MigrationFiles")
	}

	if c.DBMigrationVersion == 0 {
		return nil, errors.New("Config is missing database connection information: DBMigrationVersion")
	}

	dbName := c.DBName
	c.DBName = ""
	driver, cs := config.GetPostgresConnectionInfo(c)
	db, err := NewDB(driver, cs)

	if err != nil {
		fmt.Printf("Could not connect to database with %v %v", driver, cs)
		panic(err)
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UTC().UnixNano())
	c.DBName = fmt.Sprintf("%v_%d", dbName, rand.Int())
	CreateDatabase(db, c.DBName)
	db.Close()

	driver, cs = config.GetPostgresConnectionInfo(c)
	db, err = NewDB(driver, cs)
	if err != nil {
		fmt.Printf("Could not connect to database with %v %v", driver, cs)
		return nil, err
	}

	err = RunMigrations(db, c.DBName, c.MigrationFiles, c.DBMigrationVersion)
	if err != nil {
		fmt.Printf("Error while running migrations")
		return nil, err
	}

	return db, nil
}
