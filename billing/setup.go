package billing

import (
	"database/sql"
)

// Setup() of billing package requires:
// - Previous Setup() of auth package
// - *sql.DB's dsn has `parseTime=true`
func Setup(d *sql.DB, sqlTblPrefix string) {
	db = d
	if err := db.Ping(); err != nil {
		panic(err)
	}

	tblPrefix = sqlTblPrefix

	// Setup all tables
	setupMysqlTable()
}
