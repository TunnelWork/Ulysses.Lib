package billing

import (
	"database/sql"
	"strings"
)

/************ Shared Resources ************/

var (
	db        *sql.DB
	tblPrefix string
)

/************ Table Definitions ************/
const (
	walletTblCreation = `CREATE TABLE IF NOT EXISTS dbprefix_wallets(
        WalletID INT UNSIGNED NOT NULL AUTO_INCREMENT,
        OwnerID INT UNSIGNED DEFAULT 0, 
        Balance FLOAT NOT NULL DEFAULT 0,
        Secured FLOAT NOT NULL DEFAULT 0,
        Disabled BOOLEAN NOT NULL DEFAULT FALSE,
        PublicKey BINARY,
        PRIMARY KEY (WalletID),
        INDEX (OwnerID)
    );`
)

/************ Helper Functions ************/
func sqlStatement(query string) (*sql.Stmt, error) {
	prefixUpdatedQuery := strings.Replace(query, "dbprefix_", tblPrefix, -1)

	return db.Prepare(prefixUpdatedQuery)
}

/************ Table Creations ************/
func setupWalletTable() {
	stmt, err := sqlStatement(walletTblCreation)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
}
