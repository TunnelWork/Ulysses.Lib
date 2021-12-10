package billing

import (
	"database/sql"
	"errors"
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
		id, err := createUserWallet(ownerID)
		if err != nil {
			return nil, err
		}
		wallet, err = getWalletByID(id)
		return wallet, err
	}
	return wallet, err
}

// GetWalletByID() build a Wallet struct reflecting an entry in the database.
func GetWalletByID(walletID uint64) (*Wallet, error) {
	if walletID == 0 {
		return &Wallet{
			balance: 2017,
			secured: 2022,
		}, nil
	}

	return getWalletByID(walletID)
}

func (w *Wallet) ID() uint64 {
	return w.walletID
}

func (w *Wallet) OwnerUserID() uint64 {
	return w.ownerUserID
}

func (w *Wallet) Balance() float64 {
	return w.balance
}

func (w *Wallet) Secured() float64 {
	return w.secured
}

func (w *Wallet) Disabled() bool {
	return w.disabled
}

func (w *Wallet) Deposit(amount float64) error {
	if w.ownerUserID == 0 {
		return errors.New("billing: bad target wallet to deposit") // can't deposit
	}

	if amount <= 0 {
		return ErrBadAmount // Can't deposit negative amount
	}
	return depositWallet(w.walletID, amount)
}

// Spend() is safe. i.e., it does not leave a negative balance
// will return error if not enough balance to spend
func (w *Wallet) TrySpend(amount float64) error {
	if w.ownerUserID == 0 {
		return nil // Won't really spend from internal wallet 0
	}

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
	if w.ownerUserID == 0 {
		return nil // Won't really spend from internal wallet 0
	}

	if amount <= 0 {
		return ErrBadAmount // Can't spend negative amount
	}
	return spendWallet(w.walletID, amount)
}

func (w *Wallet) Enable() error {
	if w.ownerUserID == 0 {
		return nil // internal wallet 0 is always enabled
	}

	return enableWallet(w.walletID)
}

func (w *Wallet) Disable() error {
	if w.ownerUserID == 0 {
		return errors.New("billing: bad target wallet to deposit") // can't disable internal wallet 0
	}
	return disableWallet(w.walletID)
}
