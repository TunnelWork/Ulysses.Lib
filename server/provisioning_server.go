package server

// Server structs must be scalable, for Ulysses itself is highly scalable.
// Thus, stateless design is enforced.

// Note: Server should be self-refreshed.
// i.e., if some resource is monthly-purged,
// Server should implement required mechanisms to reflect it

type ProvisioningServer interface {
	// CreateAccount() utilizes internal server configuration and aconf pased in to create a series of accounts with
	// same sconf and variable aconf in order.
	CreateAccount(productSN uint64, accountConfiguration interface{}) error

	// ReadAccount() returns an ServerAccount-compatible struct
	GetAccount(productSN uint64) (Account, error)

	// UpdateAccount() utilizes internal server configuration and aconf pased in to update a series of accounts
	// specified by productSN.
	UpdateAccount(productSN uint64, accountConfiguration interface{}) error

	// DeleteAccount() utilizes internal server configuration to delete a series of accounts specified by accID.
	DeleteAccount(productSN uint64) error

	// Temporarily disable an account from being used or recover it
	SuspendAccount(productSN uint64) error
	UnsuspendAccount(productSN uint64) error

	// Monthly Refresh Usage
	RefreshAccount(productSN uint64) error

	// // ResourceGroup() shows usage statistics of all allocatable resources on the server
	// ResourceGroup() ServerResourceGroup
}
