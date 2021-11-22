package billing

const (
	walletTblCreation = `CREATE TABLE IF NOT EXISTS dbprefix_billing_wallets(
        wallet_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        owner_uid BIGINT UNSIGNED NOT NULL, 
        balance FLOAT NOT NULL DEFAULT 0,
        secured FLOAT NOT NULL DEFAULT 0,
        disabled BOOLEAN NOT NULL DEFAULT FALSE,
        PRIMARY KEY (wallet_id),
        INDEX (owner_uid),
        CONSTRAINT FOREIGN KEY (owner_uid) REFERENCES dbprefix_auth_user(id) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;`
)

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
func createUserWallet(ownerUserID uint64) (uint64, error) {
	stmtCreateWallet, err := sqlStatement(`INSERT INTO dbprefix_billing_wallets (owner_uid) VALUE(?);`)
	if err != nil {
		return 0, err
	}
	defer stmtCreateWallet.Close()

	result, err := stmtCreateWallet.Exec(ownerUserID)
	if err != nil {
		return 0, err
	}

	walletID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(walletID), nil
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
