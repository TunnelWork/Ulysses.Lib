package server

// ResourceID ENUMS
const (
	RESOURCE_DATA_TRANSFER uint64 = iota
	RESOURCE_SERVICE_HOUR
	RESOURCE_IP_ADDRESS_V4
	RESOURCE_IP_ADDRESS_V6
	RESOURCE_TCP_PORT
	RESOURCE_UDP_PORT
)

type Resource struct {
	ResourceID uint64 // enums should be saved somewhere in the database
	Allocated  float64
	Used       float64
	Free       float64
}
