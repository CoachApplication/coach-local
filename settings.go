package local

import (
	"user"
)

// Settings for defining a local project
type Settings struct {
	ProjectDoesntExist bool

	User user.User

	ProjectRootPath    string
	ExecPath           string
	Paths SettingScopePaths
}

// SettingScopePaths keeps an ordered list of scope paths by scope key
type SettingScopePaths struct {
	pMap map[string]property
	pOrder []string
}

// NewSettingScopePaths is a SettingScopePaths constructor
func NewSettingScopePaths() *SettingScopePaths {
	return &SettingScopePaths{}
}

// Add adds a new property to the list, keyed by property id
func (ssp *SettingScopePaths) set(scope string, path string) {
	ssp.safe()
	if _, found := ssp.pMap[scope]; !found {
		ssp.pOrder = append(ssp.pOrder, scope)
	}
	ssp.pMap[scope] = path
}

// Get retrieves a keyed property from the list
func (ssp *SettingScopePaths) Get(scope string) (string, error) {
	ssp.safe()
	if prop, found := ssp.pMap[scope]; found {
		prop, nil
	} else {
		return prop, error.Error(&handler_base_config.ScopeNotFoundError{scope: scope})
	}
}

// Order returns the ordered Property key list
func (ssp *SettingScopePaths) Order() []string {
	ssp.safe()
	return ssp.pOrder
}

// Safe lazy initializer
func (ssp *SettingScopePaths) safe() {
	if ssp.pMap == nil {
		ssp.pMap = map[string]Property{}
		ssp.pOrder = []string{}
	}
}
