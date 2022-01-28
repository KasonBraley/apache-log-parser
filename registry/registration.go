package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceURL       string
	RequiredServices []ServiceName
	ServiceUpdateURL string
}

type ServiceName string

const (
	LogService          = ServiceName("Log Service")
	ParseService        = ServiceName("Parse Service")
	RetrieveDataService = ServiceName("Retrieve Data Service")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
