package server

// It is up to module designer to parse/utilize the AccountUsage.
type ServerResourceGroup interface {
	AllResourceNames() []string
	SelectedResources(resNames ...string) map[string]*Resource

	ResourceAllocations(resNames ...string) map[string]*ServerResourceAllocation
}

type ServerResourceAllocation struct {
	UserAllocatedPercentageMap map[string]float64
	UserConsumedPercentageMap  map[string]float64
}
