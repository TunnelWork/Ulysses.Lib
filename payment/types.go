package payment

type (
	// P stands for Parameters and is a shortcut for map[string]interface{}
	P map[string]interface{}

	// PaymentRequest is an extra layer as a wrapper to provider extendability in the future
	PaymentRequest struct {
		Item PaymentUnit
	}

	PaymentResult struct {
		Status PaymentStatus
		Unit   PaymentUnit
		Msg    string
	}

	PaymentStatus uint8

	// PaymentUnit defines a single item or order to be paid
	PaymentUnit struct {
		// A caller-generated special ID used for order to track the payment
		ReferenceID string `json:"reference_id"`

		// The 3-letter currency code following ISO 4217
		// https://en.wikipedia.org/wiki/ISO_4217#Active_codes
		Currency string `json:"currency"`

		// A floating number written as a string. Precision should be limited to prevent payment issues
		Price float64 `json:"price"`
	}

	RefundRequest struct {
		Item PaymentUnit
	}
)

// OrderStatus
const (
	UNPAID PaymentStatus = iota
	PAID
	CLOSED
	UNKNOWN
)
