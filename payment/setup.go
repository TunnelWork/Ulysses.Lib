package payment

import "database/sql"

var (
	db        *sql.DB
	tblPrefix string
)

func Setup(d *sql.DB, sqlTblPrefix string) error {
	db = d
	if err := db.Ping(); err != nil {
		return err
	}

	tblPrefix = sqlTblPrefix

	return nil
}

func TblPrefix() string {
	return tblPrefix
}
