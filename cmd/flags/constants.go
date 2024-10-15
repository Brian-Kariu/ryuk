package flags

type ResourceTypes int

const (
	Workspace ResourceTypes = iota
	Environment
	Config
)

var resourceName = map[ResourceTypes]string{
	Workspace:   "workspace",
	Environment: "environment",
	Config:      "config",
}

func (rt ResourceTypes) String() string {
	return resourceName[rt]
}
