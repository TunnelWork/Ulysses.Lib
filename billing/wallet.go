package billing

import (
	"database/sql"
)

type Wallet struct {
	// Internal
	walletID    uint64
	ownerUserID uint64
	balance     float64
	secured     float64 // reserved for future use
	disabled    bool
}

func UserWallet(ownerID uint64) (*Wallet, error) {
	wallet, err := userWallet(ownerID)
	if err == sql.ErrNoRows {
		return createUserWallet(ownerID)
	}
	return wallet, err
}

// GetWalletByID() build a Wallet struct reflecting an entry in the database.
func GetWalletByID(walletID uint64) (*Wallet, error) {
	return getWalletByID(walletID)
}

func (w *Wallet) Balance() float64 {
	return w.balance
}

func (w *Wallet) Secured() float64 {
	return w.secured
}

func (w *Wallet) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrBadAmount // Can't deposit negative amount
	}
	return depositWallet(w.walletID, amount)
}

// Spend() is safe. i.e., it does not leave a negative balance
// will return error if not enough balance to spend
func (w *Wallet) TrySpend(amount float64) error {
	if amount <= 0 {
		return ErrBadAmount // Can't spend negative amount
	}
	if amount > w.balance {
		return ErrInsufficientFunds // Red Alert 2 meme
	}
	return trySpendWallet(w.walletID, amount)
}

// Spend() tries to consumes balance from the wallet, without throwing
// error even if balance is insufficient. May result in negative balance.
func (w *Wallet) Spend(amount float64) error {
	if amount <= 0 {
		return ErrBadAmount // Can't spend negative amount
	}
	return spendWallet(w.walletID, amount)
}

func (w *Wallet) Enabke() error {
	return enableWallet(w.walletID)
}

func (w *Wallet) Disable() error {
	return disableWallet(w.walletID)
}
