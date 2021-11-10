package auth

import (
	"database/sql"
)

func Setup(dbConn *sql.DB, tblPrefixOverride string) {
	db = dbConn
	tblPrefix = tblPrefixOverride
	// Create tables
	initDatabaseTable(db)
}
