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
	walletTblCreation = `CREATE TABLE IF NOT EXISTS dbprefix_billing_wallets(
        wallet_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        owner_uid BIGINT UNSIGNED NOT NULL, 
        balance FLOAT NOT NULL DEFAULT 0,
        secured FLOAT NOT NULL DEFAULT 0,
        disabled BOOLEAN NOT NULL DEFAULT FALSE,
        PRIMARY KEY (wallet_id),
        INDEX (owner_uid)
    );`

	productTblCreation = `CREATE TABLE IF NOT EXISTS dbprefix_billing_products(
        serial_number BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        owner_uid BIGINT UNSIGNED NOT NULL DEFAULT 0,
		owner_aid BIGINT UNSIGNED NOT NULL DEFAULT 0,
		product_id BIGINT UNSIGNED NOT NULL,
        date_creation DATE NOT NULL DEFAULT 0,
        date_termination DATE NOT NULL DEFAULT 0,
        billing_cycle SMALLINT UNSIGNED NOT NULL,
        wallet_id BIGINT UNSIGNED NOT NULL,
        PRIMARY KEY (serial_number),
        INDEX (owner_uid),
		INDEX (owner_aid),
		INDEX (billing_cycle)
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

/************ Product Database ************/
func addProduct(product *Product) error {
	stmt, err := sqlStatement("INSERT INTO dbprefix_billing_products(owner_uid, owner_aid, product_id, date_creation, billing_cycle, wallet_id) VALUES(?, ?, ?, CURDATE(), ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.OwnerUserID, product.OwnerAffiliationID, product.ProductID, product.BillingCycle, product.WalletID)
	return err
}

func getProductBySerialNumber(serialNumber uint64) (*Product, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_products WHERE serial_number = ?")
	if err != nil {
		return &Product{}, err
	}
	defer stmt.Close()

	var product Product
	err = stmt.QueryRow(serialNumber).Scan(&product.serialNumber, &product.OwnerUserID, &product.OwnerAffiliationID, &product.ProductID, &product.dateCreation, &product.dateTermination, &product.BillingCycle, &product.WalletID)
	return &product, err
}

func updateProduct(product *Product) error {
	stmt, err := sqlStatement("UPDATE dbprefix_billing_products SET owner_uid = ?, owner_aid = ?, product_id = ?, wallet_id = ? WHERE serial_number = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.OwnerUserID, product.OwnerAffiliationID, product.ProductID, product.WalletID, product.serialNumber)
	return err
}

func terminateProductBySerialNumber(serialNumber uint64) error {
	stmt, err := sqlStatement("UPDATE dbprefix_billing_products SET date_termination = CURDATE() WHERE serial_number = ? AND date_termination = 0")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(serialNumber)
	return err
}

func listUserProducts(ownerUserID uint64) ([]*Product, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_products WHERE owner_uid = ?")
	if err != nil {
		return []*Product{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(ownerUserID)
	if err != nil {
		return []*Product{}, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.serialNumber, &product.OwnerUserID, &product.OwnerAffiliationID, &product.ProductID, &product.dateCreation, &product.dateTermination, &product.BillingCycle, &product.WalletID)
		if err != nil {
			return []*Product{}, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func listAffiliationProducts(ownerAffiliationID uint64) ([]*Product, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_products WHERE owner_aid = ?")
	if err != nil {
		return []*Product{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(ownerAffiliationID)
	if err != nil {
		return []*Product{}, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.serialNumber, &product.OwnerUserID, &product.OwnerAffiliationID, &product.ProductID, &product.dateCreation, &product.dateTermination, &product.BillingCycle, &product.WalletID)
		if err != nil {
			return []*Product{}, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func listProductsByProductID(productID uint64) ([]*Product, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_products WHERE product_id = ?")
	if err != nil {
		return []*Product{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(productID)
	if err != nil {
		return []*Product{}, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.serialNumber, &product.OwnerUserID, &product.OwnerAffiliationID, &product.ProductID, &product.dateCreation, &product.dateTermination, &product.BillingCycle, &product.WalletID)
		if err != nil {
			return []*Product{}, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func listActiveProductsByBillingCycle(billingCycle uint8) ([]*Product, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_products WHERE billing_cycle = ? AND date_termination = 0")
	if err != nil {
		return []*Product{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(billingCycle)
	if err != nil {
		return []*Product{}, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.serialNumber, &product.OwnerUserID, &product.OwnerAffiliationID, &product.ProductID, &product.dateCreation, &product.dateTermination, &product.BillingCycle, &product.WalletID)
		if err != nil {
			return []*Product{}, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func listAllProducts() ([]*Product, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_products")
	if err != nil {
		return []*Product{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []*Product{}, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.serialNumber, &product.OwnerUserID, &product.OwnerAffiliationID, &product.ProductID, &product.dateCreation, &product.dateTermination, &product.BillingCycle, &product.WalletID)
		if err != nil {
			return []*Product{}, err
		}
		products = append(products, &product)
	}
	return products, nil
}

/************ Wallet Database ************/
func userWallet(ownerUserID uint64) (*Wallet, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_wallets WHERE owner_uid = ?")
	if err != nil {
		return &Wallet{
			disabled: true, // to prevent this bad wallet being misused
		}, err
	}
	defer stmt.Close()

	var wallet Wallet
	err = stmt.QueryRow(ownerUserID).Scan(&wallet.walletID, &wallet.ownerUserID, &wallet.balance, &wallet.secured, &wallet.disabled)
	return &wallet, err
}

// called by UserWallet() when user doesn't have a wallet
func createUserWallet(ownerUserID uint64) (*Wallet, error) {
	stmtCreateWallet, err := sqlStatement(`INSERT INTO dbprefix_billing_wallets (owner_uid) VALUE(?);`)
	if err != nil {
		return &Wallet{
			disabled: true,
		}, err
	}
	defer stmtCreateWallet.Close()

	result, err := stmtCreateWallet.Exec(ownerUserID)
	if err != nil {
		return &Wallet{
			disabled: true,
		}, err
	}

	walletID, err := result.LastInsertId()
	return &Wallet{
		walletID:    uint64(walletID),
		ownerUserID: ownerUserID,
		disabled:    false,
	}, err
}

func getWalletByID(walletID uint64) (*Wallet, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_wallets WHERE wallet_id = ?")
	if err != nil {
		return &Wallet{
			disabled: true,
		}, err
	}
	defer stmt.Close()

	var wallet Wallet
	err = stmt.QueryRow(walletID).Scan(&wallet.walletID, &wallet.ownerUserID, &wallet.balance, &wallet.secured, &wallet.disabled)
	return &wallet, err
}

func depositWallet(walletID uint64, amount float64) error {
	stmtChargeAmount, err := sqlStatement(`UPDATE dbprefix_billing_wallets 
        SET balance = balance + ? 
        WHERE wallet_id = ?;`)
	if err != nil {
		return err
	}
	defer stmtChargeAmount.Close()

	_, err = stmtChargeAmount.Exec(amount, walletID)
	return err
}

func trySpendWallet(walletID uint64, amount float64) error {
	stmtSpendAmount, err := sqlStatement(`UPDATE dbprefix_billing_wallets 
        SET balance = (CASE
            WHEN balance >= ? THEN balance - ?
            ELSE (SELECT table_name FROM information_schema.tables)
        END) WHERE wallet_id = ?;`)
	if err != nil {
		return err
	}
	defer stmtSpendAmount.Close()

	_, err = stmtSpendAmount.Exec(amount, amount, walletID)
	return err
}

func spendWallet(walletID uint64, amount float64) error {
	stmtChargeAmount, err := sqlStatement(`UPDATE dbprefix_billing_wallets 
        SET balance = balance - ? 
        WHERE wallet_id = ?;`)
	if err != nil {
		return err
	}
	defer stmtChargeAmount.Close()

	_, err = stmtChargeAmount.Exec(amount, walletID)
	return err
}

func enableWallet(walletID uint64) error {
	stmtEnable, err := sqlStatement(`UPDATE dbprefix_billing_wallets 
		SET disabled = 0 
		WHERE wallet_id = ?;`)
	if err != nil {
		return err
	}
	defer stmtEnable.Close()

	_, err = stmtEnable.Exec(walletID)
	return err
}

func disableWallet(walletID uint64) error {
	stmtDisable, err := sqlStatement(`UPDATE dbprefix_billing_wallets 
		SET disabled = 1 
		WHERE wallet_id = ?;`)
	if err != nil {
		return err
	}
	defer stmtDisable.Close()

	_, err = stmtDisable.Exec(walletID)
	return err
}
