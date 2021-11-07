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
        PublicKey TEXT,
        PRIMARY KEY (WalletID),
        INDEX (OwnerID)
    );`

	productTblCreation = `CREATE TABLE IF NOT EXISTS dbprefix_products(
        ProductID INT UNSIGNED NOT NULL AUTO_INCREMENT,
        OwnerID INT UNSIGNED NOT NULL,
        ServerType VARCHAR(32) NOT NULL,
        InstanceID VARCHAR(64) NOT NULL,
        DateCreation DATE NOT NULL DEFAULT 0,
        DateTermination DATE NOT NULL DEFAULT 0,
        BillingCycle SMALLINT NOT NULL,
        WalletID SMALLINT NOT NULL,
        PRIMARY KEY (ProductID),
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
	stmt1, err := sqlStatement(walletTblCreation)
	if err != nil {
		panic(err)
	}
	defer stmt1.Close()

	_, err = stmt1.Exec()
	if err != nil {
		panic(err)
	}

	stmt2, err := sqlStatement(productTblCreation)
	if err != nil {
		panic(err)
	}
	defer stmt2.Close()

	_, err = stmt2.Exec()
	if err != nil {
		panic(err)
	}
}
