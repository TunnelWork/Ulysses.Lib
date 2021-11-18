package server

import (
	_ "github.com/go-sql-driver/mysql"
)

type ProvisioningServerGen func( /*db *sql.DB, */ instanceID string, serverConfiguration interface{}) (ProvisioningServer, error)

var (
	psRegManagers map[string]ProvisioningServerGen = map[string]ProvisioningServerGen{}
)

// RegisterServer tells the server module how to instantiate a server struct.
// A ProvisioningServer implementation should be registered with the server module in init()
func RegisterServer(serverTypeName string, serverGen ProvisioningServerGen) {
	if psRegManagers == nil {
		psRegManagers = map[string]ProvisioningServerGen{}
	}
	psRegManagers[serverTypeName] = serverGen
}

// NewProvisioningServer returns a ProvisioningServer interface specified by serverType according to the ServerRegistrarMap
// the internal state of the returned Server interface should reflect serverConfiguration.
// A Server implementation should utilize this function to instantiate a ProvisioningServer struct with known name.
func NewProvisioningServer(serverType, instanceID string, serverConfiguration interface{}) (ProvisioningServer, error) {
	if svGen, ok := psRegManagers[serverType]; ok {
		return svGen( /*db, */ instanceID, serverConfiguration)
	} else {
		return nil, ErrServerUnknown
	}
}
