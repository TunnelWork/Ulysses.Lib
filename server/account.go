package server

type Account interface {
	Credentials() (AccountCredentials, error)
	ResourceGroup() (AccountResourceGroup, error)
}
