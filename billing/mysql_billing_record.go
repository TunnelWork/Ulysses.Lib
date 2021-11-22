package billing

const (
	billingRecordTblCreation = `CREATE TABLE IF NOT EXISTS dbprefix_billing_record(
		serial_number BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        wallet_id BIGINT UNSIGNED NOT NULL,
        user_id BIGINT UNSIGNED NOT NULL,
        product_id BIGINT UNSIGNED NOT NULL,
		product_serial_number BIGINT UNSIGNED NOT NULL,		
        billing_cycle SMALLINT UNSIGNED NOT NULL,
		billed_amount FLOAT NOT NULL,
        billed_at DATETIME NOT NULL,
        PRIMARY KEY (serial_number),
		INDEX (wallet_id),
		INDEX (user_id),
		INDEX (product_id),
		INDEX (product_serial_number),
		CONSTRAINT FOREIGN KEY (wallet_id) REFERENCES dbprefix_billing_wallets(wallet_id) ON DELETE CASCADE,
		CONSTRAINT FOREIGN KEY (user_id) REFERENCES dbprefix_auth_user(id) ON DELETE CASCADE,
		CONSTRAINT FOREIGN KEY (product_id) REFERENCES dbprefix_billing_product_listing(product_id) ON DELETE RESTRICT,
		CONSTRAINT FOREIGN KEY (product_serial_number) REFERENCES dbprefix_billing_products(serial_number) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;`
)

/**************** Billing Record Database ****************/
func addBillingRecord(record BillingRecord) (uint64, error) {
	stmt, err := sqlStatement(`INSERT INTO dbprefix_billing_record(
        wallet_id,
        user_id,
        product_id,
        product_serial_number,
        billing_cycle,
        billed_amount,
        billed_at
    ) VALUES(
        ?, ?, ?, ?, ?, ?, NOW()
    )`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		record.WalletID,
		record.UserID,
		record.ProductID,
		record.ProductSerialNumber,
		record.BillingCycle,
		record.BilledAmount,
	)
	if err != nil {
		return 0, err
	}

	serialNumber, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(serialNumber), nil
}

func listBillingRecordsByWalletID(walletID uint64) ([]BillingRecord, error) {
	stmt, err := sqlStatement(`SELECT
        serial_number,
        wallet_id,
        user_id,
        product_id,
        product_serial_number,
        billing_cycle,
        billed_amount,
        billed_at
    FROM dbprefix_billing_record
    WHERE wallet_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []BillingRecord
	for rows.Next() {
		var record BillingRecord
		err = rows.Scan(
			&record.SerialNumber,
			&record.WalletID,
			&record.UserID,
			&record.ProductID,
			&record.ProductSerialNumber,
			&record.BillingCycle,
			&record.BilledAmount,
			&record.BilledAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func listAllBillingRecords() ([]BillingRecord, error) {
	stmt, err := sqlStatement(`SELECT
        serial_number,
        wallet_id,
        user_id,
        product_id,
        product_serial_number,
        billing_cycle,
        billed_amount,
        billed_at
    FROM dbprefix_billing_record`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []BillingRecord
	for rows.Next() {
		var record BillingRecord
		err = rows.Scan(
			&record.SerialNumber,
			&record.WalletID,
			&record.UserID,
			&record.ProductID,
			&record.ProductSerialNumber,
			&record.BillingCycle,
			&record.BilledAmount,
			&record.BilledAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}
