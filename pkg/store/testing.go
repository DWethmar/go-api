package store

import (
	"database/sql"
	"fmt"

	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/database"
)

var dbCounter = 0

// WithTestStore passes a store to the provided function.
func WithTestStore(fn func(*Store)) {

	var db *sql.DB
	var myStore *Store
	con := config.LoadConfig()

	if con.DBName != "" && con.DBHost != "" && con.CreateDBScript != "" {
		dbName := con.DBName
		con.DBName = ""
		driver, cs := config.GetPostgresConnectionInfo(con)
		d, err := database.ConnectDB(driver, cs)

		if err != nil {
			fmt.Printf("Could not connect to database with %v %v", driver, cs)
			panic(err)
		}

		db = d
		dbCounter++
		con.DBName = fmt.Sprintf("%v_%v", dbName, dbCounter)

		database.CreateDatabase(db, con.DBName)
		db.Close()

		driver, cs = config.GetPostgresConnectionInfo(con)
		db, err = database.ConnectDB(driver, cs)
		if err != nil {
			fmt.Printf("Could not connect to database with %v %v", driver, cs)
			panic(err)
		}

		database.ExecSQLFile(db, con.CreateDBScript)

		fmt.Println("Using postgres repository for test server.")
		myStore = CreateStore(db)

		defer func() {
			db.Close()
			dbName := con.DBName
			con.DBName = ""
			driver, cs := config.GetPostgresConnectionInfo(con)
			db, err := database.ConnectDB(driver, cs)
			if err != nil {
				panic("Could not connect to db to drop it.")
			}
			database.DropDatabase(db, dbName)
		}()

	} else {
		fmt.Println("Using mock repository for test server.")
		myStore = CreateMockStore()
	}

	fn(myStore)
}
