package payment

import (
	"database/sql"
	"errors"
)

// A PrepaidGatewayGen is a generator function creates PrepaidGateway
// typically called NewXXXGateway() or CreateXXXGateway()
type PrepaidGatewayGen func(db *sql.DB, instanceID string, initConf interface{}) (PrepaidGateway, error)

var (
	prepaidGatewayGenRegistry map[string]PrepaidGatewayGen = map[string]PrepaidGatewayGen{}
	prepaidGatewayRegistry    map[string]PrepaidGateway    = map[string]PrepaidGateway{}
)

func RegisterPrepaidGatewayGenerator(gatewayTypeName string, genFunc PrepaidGatewayGen) {
	// if _, ok := prepaidGatewayRegistry[gatewayTypeName]; ok {
	// 	panic(fmt.Sprintf("payment.RegisterPrepaidGateway() tries to register repeated gatewayTypeName: %s", gatewayTypeName)) // This is a panic() level conflict
	// }
	prepaidGatewayGenRegistry[gatewayTypeName] = genFunc // Overwrites existing ones
}

// NewPrepaidGateway() creates a PrepaidGateway in the type of gatewayTypeName
// with the specified instanceID, which supports the duplicapability for each gatewayType.
// Note: it is caller's responsibility to make sure the *sql.DB is alive.
func NewPrepaidGateway(gatewayTypeName, instanceID string, initConf interface{}) (PrepaidGateway, error) {
	if genFunc, ok := prepaidGatewayGenRegistry[gatewayTypeName]; ok {
		gateway, err := genFunc(db, instanceID, initConf)
		if err != nil {
			return nil, err
		}
		prepaidGatewayRegistry[instanceID] = gateway
		return gateway, nil
	} else {
		return nil, errors.New("payment: gateway name not found")
	}
}

func RegisterPrepaidGateway(instanceID string, gateway PrepaidGateway) {
	prepaidGatewayRegistry[instanceID] = gateway
}

func GetPrepaidGateway(instanceID string) (PrepaidGateway, error) {
	if gateway, ok := prepaidGatewayRegistry[instanceID]; ok {
		return gateway, nil
	} else {
		return nil, errors.New("payment: gateway not found")
	}
}
