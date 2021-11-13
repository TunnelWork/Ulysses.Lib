package server

// It is up to module designer to parse/utilize the Credential.
type Credentials interface {
	// Customer() returns the customer-oriented information.
	// interface{} types should be in string/number/bool/[]string/[]number
	Customer() []*Credential

	// Provider() returns the provider-oriented information.
	// interface{} types should be in string/number/bool/[]string/[]number
	Admin() []*Credential
}
