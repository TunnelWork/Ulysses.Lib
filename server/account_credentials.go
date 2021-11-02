package server

// It is up to module designer to parse/utilize the Credential.
type AccountCredentials interface {
	Customer() (credentials map[string]string)
	Admin() (credentials map[string]string)
}
