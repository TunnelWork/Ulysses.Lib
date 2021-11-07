package billing

import (
	"database/sql"
	"encoding/base64"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type Wallet struct {
	// Internal
	walletID uint64
	ownerID  uint64

	publicKey []byte
}

func UserWallet(ownerID uint64) (*Wallet, error) {
	stmtCheckoutWallet, err := sqlStatement(`SELECT WalletID, PublicKey FROM dbprefix_wallets WHERE OwnerID = ? AND Disabled = FALSE ORDER BY WalletID ASC;`)
	if err != nil {
		return nil, err
	}
	defer stmtCheckoutWallet.Close()

	var walletID uint64
	var publicKeyB64 sql.NullString
	err = stmtCheckoutWallet.QueryRow(ownerID).Scan(&walletID, &publicKeyB64)
	if err != sql.ErrNoRows {
		resultWallet := Wallet{
			walletID: walletID,
			ownerID:  ownerID,
		}
		if publicKeyB64.Valid {
			publicKey, _ := base64.StdEncoding.DecodeString(publicKeyB64.String)
			resultWallet.publicKey = publicKey
		}
		return &resultWallet, err
	} else { // register new wallet for user
		return createWallet(ownerID)
	}
}

// called by UserWallet() when user doesn't have a wallet
func createWallet(ownerID uint64) (*Wallet, error) {
	stmtCreateWallet, err := sqlStatement(`INSERT INTO dbprefix_wallets (OwnerID) VALUE(?);`)
	if err != nil {
		return nil, err
	}
	defer stmtCreateWallet.Close()

	result, err := stmtCreateWallet.Exec(ownerID)
	if err != nil {
		return nil, err
	}

	walletID, err := result.LastInsertId()
	return &Wallet{
		walletID: uint64(walletID),
		ownerID:  ownerID,
	}, err
}

// WalletByID() build a Wallet struct reflecting an entry in the database.
func WalletByID(walletID uint64) (*Wallet, error) {
	stmtCheckoutWallet, err := sqlStatement(`SELECT OwnerID, PublicKey FROM dbprefix_wallets WHERE WalletID = ? AND Disabled = FALSE;`)
	if err != nil {
		return nil, err
	}
	defer stmtCheckoutWallet.Close()

	var ownerID uint64
	var publicKeyB64 sql.NullString
	err = stmtCheckoutWallet.QueryRow(walletID).Scan(&ownerID, &publicKeyB64)
	if err != sql.ErrNoRows {
		resultWallet := Wallet{
			walletID: walletID,
			ownerID:  ownerID,
		}
		if publicKeyB64.Valid {
			publicKey, _ := base64.StdEncoding.DecodeString(publicKeyB64.String)
			resultWallet.publicKey = publicKey
		}
		return &resultWallet, err
	} else { // register new wallet for user
		return nil, errors.New("billing: wallet not found")
	}
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

// Use Ed25519
// TODO: Upload to Database
func (w *Wallet) UploadPubkey(pubkey []byte, signature string /*, msg string = "UploadPubkey"*/) error {
	var msg string = "UploadPubkey"

	ed25519Master := jwt.SigningMethodEd25519{}
	err := ed25519Master.Verify(msg, signature, pubkey)
	return err
}

// TODO: VerifySignature()
func (w *Wallet) VerifySignature(signature, data string) error {
	return nil
}
