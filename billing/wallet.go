package billing

type Wallet struct {
	// Internal
	walletID uint64
	ownerID  uint64

	// Total assets owned by a user
	// Assets secured for an active Pay-As-You-Go service
}

func UserWallet(ownerID uint64) (*Wallet, error) {
	stmtCheckoutWallet, err := sqlStatement(`SELECT WalletID FROM dbprefix_wallets WHERE OwnerID = ? AND Disabled = FALSE ORDER BY WalletID ASC;`)
	if err != nil {
		return nil, err
	}
	defer stmtCheckoutWallet.Close()

	var walletID uint64
	err = stmtCheckoutWallet.QueryRow(ownerID).Scan(&walletID)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		walletID: walletID,
		ownerID:  ownerID,
	}, nil
}

func (w *Wallet) Balance() (float64, error) {
	stmtRealtimeBalance, err := sqlStatement(`SELECT Balance FROM dbprefix_wallets where WalletID = ?;`)
	if err != nil {
		return 0, err
	}
	defer stmtRealtimeBalance.Close()

	var balance float64
	err = stmtRealtimeBalance.QueryRow(w.walletID).Scan(&balance)
	return balance, err
}

func (w *Wallet) Secured() (float64, error) {
	stmtRealtimeSecuredFund, err := sqlStatement(`SELECT Secured FROM dbprefix_wallets where WalletID = ?;`)
	if err != nil {
		return 0, err
	}
	defer stmtRealtimeSecuredFund.Close()

	var secured float64
	err = stmtRealtimeSecuredFund.QueryRow(w.walletID).Scan(&secured)
	return secured, err
}

func (w *Wallet) AvailableFund() (float64, error) {
	stmtRealtimeAvailableFund, err := sqlStatement(`SELECT Balance-Secured FROM dbprefix_wallets where WalletID = ?;`)
	if err != nil {
		return 0, err
	}
	defer stmtRealtimeAvailableFund.Close()

	var available float64
	err = stmtRealtimeAvailableFund.QueryRow(w.walletID).Scan(&available)
	return available, err
}

// Spend() is safe. i.e., it does not leave a negative balance
// will return error if not enough balance to spend
func (w *Wallet) Spend(amount float64) error {
	if amount <= 0 {
		return ErrBadAmount
	}
	stmtSpendAmount, err := sqlStatement(`UPDATE dbprefix_wallets 
        SET Balance = (CASE
            WHEN Balance >= ? THEN Balance - ?
            ELSE (SELECT table_name FROM information_schema.tables)
        END) WHERE WalletID = ?;`)
	if err != nil {
		return err
	}
	defer stmtSpendAmount.Close()

	_, err = stmtSpendAmount.Exec(amount, amount, w.walletID)
	return err
}

// Charge() is unsafe version of Spend().
// It doesn't fail when user can't afford it.
// USE ONLY FOR POST-USE BILLING
func (w *Wallet) Charge(amount float64) error {
	if amount <= 0 {
		return ErrBadAmount
	}
	stmtChargeAmount, err := sqlStatement(`UPDATE dbprefix_wallets 
        SET Balance = Balance - ? 
        WHERE WalletID = ?;`)
	if err != nil {
		return err
	}
	defer stmtChargeAmount.Close()

	_, err = stmtChargeAmount.Exec(amount, w.walletID)
	return err
}

func (w *Wallet) SecureFund(amount float64) error {
	if amount <= 0 {
		return ErrBadAmount
	}
	stmtSecureAmount, err := sqlStatement(`UPDATE dbprefix_wallets 
        SET Balance = (CASE
            WHEN Balance >= ? THEN Balance - ?
            ELSE (SELECT table_name FROM information_schema.tables)
        END), Secured = Secured + ? WHERE WalletID = ?;`)
	if err != nil {
		return err
	}
	defer stmtSecureAmount.Close()

	_, err = stmtSecureAmount.Exec(amount, amount, amount, w.walletID)
	return err
}

func (w *Wallet) UndoSecureFund(amount float64) error {
	if amount <= 0 {
		return ErrBadAmount
	}
	stmtSecureAmount, err := sqlStatement(`UPDATE dbprefix_wallets 
        SET Secured = (CASE
            WHEN Secured >= ? THEN Secured - ?
            ELSE (SELECT table_name FROM information_schema.tables)
        END), Balance = Balance + ? WHERE WalletID = ?;`)
	if err != nil {
		return err
	}
	defer stmtSecureAmount.Close()

	_, err = stmtSecureAmount.Exec(amount, amount, amount, w.walletID)
	return err
}

func (w *Wallet) SpendSecured(amount float64) error {
	if amount <= 0 {
		return ErrBadAmount
	}
	stmtSpendSecuredAmount, err := sqlStatement(`UPDATE dbprefix_wallets 
        SET Secured = (CASE
            WHEN Secured >= ? THEN Secured - ?
            ELSE (SELECT table_name FROM information_schema.tables)
        END) WHERE WalletID = ?;`)
	if err != nil {
		return err
	}
	defer stmtSpendSecuredAmount.Close()

	_, err = stmtSpendSecuredAmount.Exec(amount, amount, w.walletID)
	return err
}

func (w *Wallet) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrBadAmount
	}
	stmtChargeAmount, err := sqlStatement(`UPDATE dbprefix_wallets 
        SET Balance = Balance + ? 
        WHERE WalletID = ?;`)
	if err != nil {
		return err
	}
	defer stmtChargeAmount.Close()

	_, err = stmtChargeAmount.Exec(amount, w.walletID)
	return err
}
