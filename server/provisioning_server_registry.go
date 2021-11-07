package server

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type ProvisioningServerGen func(db *sql.DB, instanceID string, serverConfiguration interface{}) (ProvisioningServer, error)

var (
	psRegManagers map[string]ProvisioningServerGen = map[string]ProvisioningServerGen{}
)

// RegisterServer tells the server module how to instantiate a server struct
func RegisterServer(serverTypeName string, serverGen ProvisioningServerGen) {
	if psRegManagers == nil {
		psRegManagers = map[string]ProvisioningServerGen{}
	}
	psRegManagers[serverTypeName] = serverGen
}

// NewServer returns a Server interface specified by serverType according to the ServerRegistrarMap
// the internal state of the returned Server interface should reflect sconf.
func NewServer(db *sql.DB, serverType string, instanceID string, serverConfiguration interface{}) (ProvisioningServer, error) {
	if svGen, ok := psRegManagers[serverType]; ok {
		return svGen(db, instanceID, serverConfiguration)
	} else {
		return nil, ErrServerUnknown
	}
}
