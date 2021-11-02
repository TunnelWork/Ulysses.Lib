package server

import "database/sql"

type ServerGen func(db *sql.DB, instanceID string, serverConfiguration interface{}) (Server, error)

var (
	regManagers map[string]ServerGen = map[string]ServerGen{}
)

// RegisterServer tells the server module how to instantiate a server struct
func RegisterServer(serverTypeName string, serverGen ServerGen) {
	if regManagers == nil {
		regManagers = map[string]ServerGen{}
	}
	regManagers[serverTypeName] = serverGen
}

// NewServer returns a Server interface specified by serverType according to the ServerRegistrarMap
// the internal state of the returned Server interface should reflect sconf.
func NewServer(db *sql.DB, serverType string, instanceID string, serverConfiguration interface{}) (Server, error) {
	if svGen, ok := regManagers[serverType]; ok {
		return svGen(db, instanceID, serverConfiguration)
	} else {
		return nil, ErrServerUnknown
	}
}
