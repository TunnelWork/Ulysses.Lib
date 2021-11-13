package billing

import "errors"

var (
	ErrBadAmount         error = errors.New("billing: bad amount input")
	ErrInsufficientFunds error = errors.New("billing: insufficient funds") // I can hear it...
)
