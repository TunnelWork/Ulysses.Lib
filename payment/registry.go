package payment

import (
	"database/sql"
	"errors"
)

// A PrepaidGatewayGen is a generator function creates PrepaidGateway
// typically called NewXXXGateway() or CreateXXXGateway()
type PrepaidGatewayGen func(db *sql.DB, instanceID string, initConf interface{}) (PrepaidGateway, error)

var (
	prepaidGatewayRegistry map[string]PrepaidGatewayGen = map[string]PrepaidGatewayGen{}
)

func RegisterPrepaidGateway(gatewayTypeName string, genFunc PrepaidGatewayGen) {
	// if _, ok := prepaidGatewayRegistry[gatewayTypeName]; ok {
	// 	panic(fmt.Sprintf("payment.RegisterPrepaidGateway() tries to register repeated gatewayTypeName: %s", gatewayTypeName)) // This is a panic() level conflict
	// }
	prepaidGatewayRegistry[gatewayTypeName] = genFunc // Overwrites existing ones
}

// NewPrepaidGateway() creates a PrepaidGateway in the type of gatewayTypeName
// with the specified instanceID, which supports the duplicapability for each gatewayType.
// Note: it is caller's responsibility to make sure the *sql.DB is alive.
func NewPrepaidGateway(db *sql.DB, gatewayTypeName string, instanceID string, initConf interface{}) (PrepaidGateway, error) {
	if genFunc, ok := prepaidGatewayRegistry[gatewayTypeName]; ok {
		return genFunc(db, instanceID, initConf)
	} else {
		return nil, errors.New("payment: gateway name not found")
	}
}
