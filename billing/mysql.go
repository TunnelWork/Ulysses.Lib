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

/************ Helper Functions ************/
func sqlStatement(query string) (*sql.Stmt, error) {
	prefixUpdatedQuery := strings.ReplaceAll(query, "dbprefix_", tblPrefix)

	return db.Prepare(prefixUpdatedQuery)
}

/************ Table Creations ************/
func setupMysqlTable() {
	// dbprefix_billing_wallets relys on dbprefix_auth_user
	stmt1, err := sqlStatement(walletTblCreation)
	if err != nil {
		panic(err)
	}
	defer stmt1.Close()

	_, err = stmt1.Exec()
	if err != nil {
		panic(err)
	}

	// dbprefix_billing_product_listing_group relys on no table
	stmt2, err := sqlStatement(productListingGroupTblCreation)
	if err != nil {
		panic(err)
	}
	defer stmt2.Close()

	_, err = stmt2.Exec()
	if err != nil {
		panic(err)
	}

	// dbprefix_billing_product_listing relys on dbprefix_billing_product_listing_group
	stmt3, err := sqlStatement(productListingTblCreation)
	if err != nil {
		panic(err)
	}
	defer stmt3.Close()

	_, err = stmt3.Exec()
	if err != nil {
		panic(err)
	}

	// dbprefix_billing_products relys on:
	// - dbprefix_billing_product_listing
	// - dbprefix_auth_user
	// - dbprefix_auth_affiliation
	stmt4, err := sqlStatement(productTblCreation)
	if err != nil {
		panic(err)
	}
	defer stmt4.Close()

	_, err = stmt4.Exec()
	if err != nil {
		panic(err)
	}

	// dbprefix_billing_record relys on:
	// - dbprefix_auth_user
	// - dbprefix_billing_wallets
	// - dbprefix_billing_product_listing
	// - dbprefix_billing_products
	stmt5, err := sqlStatement(billingRecordTblCreation)
	if err != nil {
		panic(err)
	}
	defer stmt5.Close()

	_, err = stmt5.Exec()
	if err != nil {
		panic(err)
	}
}
