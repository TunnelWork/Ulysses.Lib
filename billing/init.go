package billing

import "database/sql"

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
