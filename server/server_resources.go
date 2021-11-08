package server

// It is up to module designer to parse/utilize the AccountUsage.
type ServerResourceGroup interface {
	ListResources() []uint64                              // uint64 is enum of resource name
	SelectedResources(res ...uint64) map[uint64]*Resource // uint64 is enum of resource name

	ResourceAllocations(res ...uint64) map[uint64]*ServerResourceAllocation //	uint64 is enum of resource name
}

type ServerResourceAllocation struct {
	ProductAllocatedPercentageMap map[uint64]float64 // uint64 is product id
	ProductConsumedPercentageMap  map[uint64]float64 // uint64 is product id
}
