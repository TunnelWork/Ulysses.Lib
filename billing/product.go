package billing

type Product struct {
	Pid uint64

	// Ownership
	OwnerID uint64

	// Server-Account binding
	ServerType string
	InstanceID string
	// referenceID uint64 // use Pid instead!

	// billing/payment
	BillingCycle BillingCycle

	//// Only for BillingCycle: PayAsYouGo
	HourlyRate              float64
	MonthlySpendingCap      float64
	CurrentMonthExpenditure float64
}

type BillingCycle uint8

const (
	PayAsYouGo BillingCycle = iota
	PerMonth
	PerQuarter
	PerYear
)
