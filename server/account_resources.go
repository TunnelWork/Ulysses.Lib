package server

// It is up to module designer to parse/utilize the AccountUsage.
type AccountResourceGroup interface {
	ListResources() []uint64
	SelectedResources(res ...uint64) map[uint64]*Resource
}
