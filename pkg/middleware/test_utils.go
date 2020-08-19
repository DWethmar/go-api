package middleware

import (
	"database/sql"
	"fmt"

	"github.com/dwethmar/go-api/pkg/config"
	"github.com/dwethmar/go-api/pkg/database"
	"github.com/dwethmar/go-api/pkg/store"
)

var dbCounter = 0

func withStore(fn func(*store.Store)) {

	var db *sql.DB
	var myStore *store.Store
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

		database.ExecSQLFileDatabase(db, con.CreateDBScript)

		fmt.Println("Using postgres repository for test server.")
		myStore = store.CreateStore(db)

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
		myStore = store.CreateMockStore()
	}

	fn(myStore)
}
