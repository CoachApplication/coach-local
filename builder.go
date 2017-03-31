package local

import (
	"context"

	api "github.com/CoachApplication/api"
	base "github.com/CoachApplication/base"
	config_provider "github.com/CoachApplication/config/provider"
)

// Builder Standard local coach api.Builder
type Builder struct {
	context  context.Context
	settings Settings
	parent   api.API

	implementations []string

	sharedConfigProvider config_provider.Provider
}

// NewBuilder Constructor for Builder from Settings
func NewBuilder(ctx context.Context, settings Settings) *Builder {
	return &Builder{
		context:         ctx,
		settings:        settings,
		implementations: []string{},
	}
}

// Builder explicilty convert this struct to an api.Builder
func (b *Builder) Builder() api.Builder {
	return api.Builder(b)
}

// Id provides a unique machine name for the Builder
func (b *Builder) Id() string {
	return "local.standard"
}

// SetParent Provides the API reference to the Builder which may use it's operations internally
func (b *Builder) SetParent(parent api.API) {
	b.parent = parent
}

// Activate Enable keyed implementations, providing settings for those handler implementations
func (b *Builder) Activate(implementations []string, settings api.SettingsProvider) error {
	for _, implementation := range implementations {
		found := false
		for _, exist := range b.implementations {
			if exist == implementation {
				found = true
				break
			}
		}
		if !found {
			b.implementations = append(implementations, implementation)
		}
	}

	return nil
}

// Validates Ask the builder if it is happy and willing to provide operations
func (b *Builder) Validate() api.Result {
	return base.MakeSuccessfulResult()
}

// Operations provide any Builder user with a set of Operation objects
func (b *Builder) Operations() api.Operations {
	ops := base.NewOperations()

	for _, imp := range b.implementations {
		switch imp {
		case "config":
			ops.Merge(b.configOperations())
		case "orchestrate":
			ops.Merge(b.orchestrateOperations())
		case "project":
			ops.Merge(b.projectOperations())
		}
	}

	return ops.Operations()
}
