package billing

import "database/sql"

const (
	productTblCreation = `CREATE TABLE IF NOT EXISTS dbprefix_billing_products(
        serial_number BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        owner_uid BIGINT UNSIGNED NOT NULL DEFAULT 0,
        owner_aid BIGINT UNSIGNED NOT NULL DEFAULT 0,
        product_id BIGINT UNSIGNED NOT NULL,
        date_creation DATE NOT NULL DEFAULT 0,
        date_last_bill DATE NOT NULL DEFAULT 0,
        date_termination DATE NOT NULL DEFAULT 0,
        already_terminated BOOLEAN NOT NULL DEFAULT FALSE,
        wallet_id BIGINT UNSIGNED NOT NULL,
        billing_cycle SMALLINT UNSIGNED NOT NULL,
        price FLOAT NOT NULL,
        monthly_spending_cap FLOAT NOT NULL,
        current_month_spending FLOAT NOT NULL,
        PRIMARY KEY (serial_number),
        INDEX (product_id),
        INDEX (owner_uid),
        INDEX (owner_aid),
        INDEX (billing_cycle),
        CONSTRAINT FOREIGN KEY (product_id) REFERENCES dbprefix_billing_product_listing(product_id) ON DELETE RESTRICT,
        CONSTRAINT FOREIGN KEY (owner_uid) REFERENCES dbprefix_auth_user(id) ON DELETE CASCADE,
        CONSTRAINT FOREIGN KEY (owner_aid) REFERENCES dbprefix_auth_affiliation(id) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;`
)

/************ Product Database ************/
func rowToProduct(row *sql.Row) (*Product, error) {
	var product Product
	err := row.Scan(
		&product.serialNumber,
		&product.OwnerUserID,
		&product.OwnerAffiliationID,
		&product.ProductID,
		&product.dateCreation,
		&product.dateLastBill,
		&product.dateTermination,
		&product.terminated,
		&product.WalletID,
		&product.BillingOption.BillingCycle,
		&product.BillingOption.Price,
		&product.BillingOption.MonthlySpendingCap,
		&product.BillingOption.CurrentMonthSpending,
	)
	return &product, err
}

func rowsToProductSlice(rows *sql.Rows) ([]*Product, error) {
	var products []*Product
	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.serialNumber,
			&product.OwnerUserID,
			&product.OwnerAffiliationID,
			&product.ProductID,
			&product.dateCreation,
			&product.dateLastBill,
			&product.dateTermination,
			&product.terminated,
			&product.WalletID,
			&product.BillingOption.BillingCycle,
			&product.BillingOption.Price,
			&product.BillingOption.MonthlySpendingCap,
			&product.BillingOption.CurrentMonthSpending,
		)
		if err != nil {
			return []*Product{}, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func addProduct(product *Product) (serialNumber uint64, err error) {
	stmt, err := sqlStatement(`INSERT INTO 
    dbprefix_billing_products(
        owner_uid, owner_aid, product_id, date_creation, date_last_bill, wallet_id, billing_cycle, price, monthly_spending_cap
    ) VALUES(?, ?, ?, CURDATE(), CURDATE(), ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		product.OwnerUserID,
		product.OwnerAffiliationID,
		product.ProductID,
		product.WalletID,
		product.BillingOption.BillingCycle,
		product.BillingOption.Price,
		product.BillingOption.MonthlySpendingCap,
		// product.BillingOption.CurrentMonthSpending, // 0 for new products
	)

	if err != nil {
		return 0, err
	}

	serialNumberInt64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	serialNumber = uint64(serialNumberInt64)

	return serialNumber, nil
}

func getProductBySerialNumber(serialNumber uint64) (*Product, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_products WHERE serial_number = ?")
	if err != nil {
		return &Product{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(serialNumber)
	product, err := rowToProduct(row)
	return product, err
}

func updateProduct(product *Product) error {
	stmt, err := sqlStatement(`UPDATE dbprefix_billing_products SET 
    owner_uid = ?, 
    owner_aid = ?, 
    product_id = ?, 
    wallet_id = ?,    
    price = ?, 
    current_month_spending = ? 
    WHERE serial_number = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		product.OwnerUserID,
		product.OwnerAffiliationID,
		product.ProductID,
		product.WalletID,
		product.BillingOption.Price,
		product.BillingOption.CurrentMonthSpending,
		product.serialNumber,
	)
	return err
}

func terminateProductBySerialNumber(serialNumber uint64) error {
	stmt, err := sqlStatement("UPDATE dbprefix_billing_products SET date_termination = CURDATE() WHERE serial_number = ? AND already_terminated = FALSE")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(serialNumber)
	return err
}

func toTerminateProductOn(serialNumber uint64, date string) error {
	stmt, err := sqlStatement("UPDATE dbprefix_billing_products SET date_termination = ? WHERE serial_number = ? AND already_terminated = FALSE")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(date, serialNumber)
	return err
}

func updateDateLastBillProductBySerialNumber(serialNumber uint64) error {
	stmt, err := sqlStatement("UPDATE dbprefix_billing_products SET date_last_bill = CURDATE() WHERE serial_number = ?")
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

	products, err := rowsToProductSlice(rows)
	return products, err
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

	products, err := rowsToProductSlice(rows)
	return products, err
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

	products, err := rowsToProductSlice(rows)
	return products, err
}

func listActiveProductsByBillingCycle(billingCycle uint8) ([]*Product, error) {
	// List all products that is either not to be terminated yet (date_termination = 0)
	// or has not been terminated yet (date_termination > today)
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_products WHERE billing_cycle = ? AND (date_termination = 0 OR date_termination > CURDATE()) AND already_terminated = FALSE")
	if err != nil {
		return []*Product{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(billingCycle)
	if err != nil {
		return []*Product{}, err
	}
	defer rows.Close()

	products, err := rowsToProductSlice(rows)
	return products, err
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

	products, err := rowsToProductSlice(rows)
	return products, err
}

func listProductsToTerminate() ([]*Product, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_products WHERE (date_termination < CURDATE() OR date_termination = CURDATE()) AND already_terminated = FALSE")
	if err != nil {
		return []*Product{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []*Product{}, err
	}
	defer rows.Close()

	products, err := rowsToProductSlice(rows)
	return products, err
}
