package server

type ServerAccount interface {
	Credentials() (AccountCredentials, error)
	ResourceGroup() (AccountResourceGroup, error)
}
