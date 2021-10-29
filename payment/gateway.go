package payment

// Note: Gateway implementations should be scalable, as Ulysses itself is highly scalable.
// It is recommended to save anything possible to database.

// PrepaidGateway allows user to pay
// based on their purchase or deposit order.
type PrepaidGateway interface {
	/**** Pay ****/
	CheckoutForm(pr PaymentRequest) (HTMLCheckoutForm string, err error)

	/**** Status ****/
	// PaymentStatus() checks for a referenceID
	// this function should be called once a customer CLAIMS the payment has been made
	PaymentStatus(referenceID string) (paymentStatus PaymentStatus, err error)

	/**** Refund ****/
	// IsRefundable() checks for refundability for a referenceID
	// It should always return false for a gateway without Refund() capability
	IsRefundable(referenceID string) bool

	// Refund() returns nil if successfully refunded.
	Refund(rr RefundRequest) error

	/**** Callback ****/
	// OnStatusChange() sets the function to be called once the referenceID's payment status is changed
	// returns error when doesn't have such callback functionality
	OnStatusChange(SuccessFunc func(referenceID string, newPaymentStatus PaymentStatus)) error
}
