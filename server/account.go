package server

type Account interface {
	Credentials() (Credentials, error)
	Resources() ([]*Resource, error)
}
