package registry

type Registration struct {
	ServiceName ServiceName
	ServiceURL  string
}

type ServiceName string

const (
	LogService          = ServiceName("Log Service")
	ParseService        = ServiceName("Parse Service")
	RetrieveDataService = ServiceName("Retrieve Data Service")
)
