package billing

import "errors"

var (
	ErrBadAmount error = errors.New("billing: bad amount input")
)
