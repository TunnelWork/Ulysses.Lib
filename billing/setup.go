package billing

import (
	"database/sql"
)

// Setup() of billing package requires:
// - Previous Setup() of auth package
// - *sql.DB's dsn has `parseTime=true`
func Setup(d *sql.DB, sqlTblPrefix string) error {
	if err := db.Ping(); err != nil {
		return err
	}

	db = d
	tblPrefix = sqlTblPrefix

	// Setup all tables
	setupWalletTable()

	return nil
}
