package server

// It is up to module designer to parse/utilize the AccountUsage.
type AccountResourceGroup interface {
	AllResourceNames() []string
	SelectedResources(resNames ...string) map[string]*Resource
}
