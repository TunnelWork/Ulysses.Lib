package auth

import (
	"database/sql"
)

func Setup(dbConn *sql.DB, tblPrefixOverride string) {
	db = dbConn
	if db.Ping() != nil {
		panic("Could not connect to database")
	}

	tblPrefix = tblPrefixOverride
	// Create tables
	initDatabaseTable(db)
}
