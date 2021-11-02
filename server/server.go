package server

// Server structs must be scalable, for Ulysses itself is highly scalable.
// Thus, stateless design is enforced.

// Note: Server should be self-refreshed.
// i.e., if some resource is monthly-purged,
// Server should implement required mechanisms to reflect it

type Server interface {
	/// Server

	// ResourceGroup() shows usage statistics of all allocatable resources on the server
	ResourceGroup() ServerResourceGroup

	//// Account

	// CreateAccount() utilizes internal server configuration and aconf pased in to create a series of accounts with
	// same sconf and variable aconf in order.
	CreateAccount(referenceID int, accountConfiguration interface{}) error

	// ReadAccount() returns an ServerAccount-compatible struct
	ReadAccount(referenceID int) (ServerAccount, error)

	// UpdateAccount() utilizes internal server configuration and aconf pased in to update a series of accounts
	// specified by referenceID.
	UpdateAccount(referenceID int, accountConfiguration interface{}) error

	// DeleteAccount() utilizes internal server configuration to delete a series of accounts specified by accID.
	DeleteAccount(referenceID int) error

	// PurgeResourceUsage sets all USED resource counter to 0 for all users.
	// usecase: clean-reinstall
	// Won't be automatically called on a time-basis. Not a cronjob mounting point.
	PurgeResourceUsage()
}
