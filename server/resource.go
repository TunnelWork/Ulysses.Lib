package server

// ResourceID ENUMS
const (
	RESOURCE_DATA_TRANSFER uint64 = iota + 1
	RESOURCE_SERVICE_HOUR
	RESOURCE_IP_ADDRESS_V4
	RESOURCE_IP_ADDRESS_V6
	RESOURCE_TCP_PORT
	RESOURCE_UDP_PORT
)

type Resource struct {
	ResourceID uint64  `json:"resource_id"` // enums should be saved somewhere in the database
	Allocated  float64 `json:"allocated"`
	Used       float64 `json:"used"`
	Free       float64 `json:"free"`
}
