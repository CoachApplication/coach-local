package local

import (
	coach_api "github.com/james-nesbitt/coach-api"
	coach_base "github.com/james-nesbitt/coach-base"
)

// Builder Standard local coach api.Builder
type Builder struct {
	settings Settings
	parent   coach_api.API

	implementations []string
}

// NewBuilder Constructor for Builder from Settings
func NewBuilder(settings Settings) *Builder {
	return &Builder{
		settings:        settings,
		implementations: []string{},
	}
}

// Builder explicilty convert this struct to an api.Builder
func (b *Builder) Builder() coach_api.Builder {
	return coach_api.Builder(b)
}

// Id provides a unique machine name for the Builder
func (b *Builder) Id() string {
	return "local.standard"
}

// SetParent Provides the API reference to the Builder which may use it's operations internally
func (b *Builder) SetParent(parent coach_api.API) {
	b.parent = parent
}

// Activate Enable keyed implementations, providing settings for those handler implementations
func (b *Builder) Activate(implementations []string, settings coach_api.SettingsProvider) error {
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
func (b *Builder) Validate() coach_api.Result {
	return coach_base.MakeSuccessfulResult()
}

// Operations provide any Builder user with a set of Operation objects
func (b *Builder) Operations() coach_api.Operations {
	ops := coach_base.NewOperations()

	return ops.Operations()
}

func (b *Builder) configOperations() {
	configOperations := handler_configconnector.GetOperations(
		handler_configconnect.NewfileConfigConnector(
			b.settings.Path.Order()
			func(path, scope string) string { return path.Join(path, scope + ".yml")},
		),
		handler_configconnect.NewymlConfigInterpreter(

		),
	)
}