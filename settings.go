package local

import (
	"errors"
	"os/user"
)

// Settings for defining a local project
type Settings struct {
	ProjectDoesntExist bool

	User user.User

	ProjectRootPath string
	ExecPath        string
	Paths           SettingScopePaths
}

// SettingScopePaths keeps an ordered list of scope paths by scope key
type SettingScopePaths struct {
	pMap   map[string]string
	pOrder []string
}

// NewSettingScopePaths is a SettingScopePaths constructor
func NewSettingScopePaths() *SettingScopePaths {
	return &SettingScopePaths{}
}

// Add adds a new property to the list, keyed by property id
func (ssp *SettingScopePaths) Set(scope string, path string) {
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
		return prop, nil
	} else {
		return prop, errors.New("Scope not found")
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
		ssp.pMap = map[string]string{}
		ssp.pOrder = []string{}
	}
}
