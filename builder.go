package local

import (
	api "github.com/james-nesbitt/coach-api"
	base "github.com/james-nesbitt/coach-base"
	base_config_provider "github.com/james-nesbitt/coach-config/provider"
	base_config_provider_file "github.com/james-nesbitt/coach-config/provider/file"
	base_config_provider_yaml "github.com/james-nesbitt/coach-config/provider/yaml"
)

// Builder Standard local coach api.Builder
type Builder struct {
	settings Settings
	parent   api.API

	implementations []string

	sharedConfigProvider base_config_provider.Provider
}

// NewBuilder Constructor for Builder from Settings
func NewBuilder(settings Settings) *Builder {
	return &Builder{
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

	return ops.Operations()
}

func (b *Builder) configProvider() base_config_provider.Provider {
	if b.sharedConfigProvider == nil {
		/**
		 * Stay with me here
		 *
		 * Config here is going to come from a Backend multiplexing config provider, but it is going to use
		 * only a single backend, which will be a FileConnector, that loads YML files based on an arrangement
		 * of files where subfolders are used for different scopes, and the YML files are named after the
		 * config keys.
		 * In the local case, each of the settings Paths is considered a valid scope.
		 *
		 * Why so much complexity?
		 *
		 * Because it should be easy to layer backends so that we can provide defaults for some cases, like
		 * missing files for some config.  We will want to add a Default provider very soon.
		 * There are many options for additional backend, included buffered connectors, but there is not yet
		 * any fail-over mechanism for a backend.  This is something that is needed.
		 */

		// This provider
		tbcp := base_config_provider.NewBackendConfigProvider()
		b.sharedConfigProvider = tbcp.Provider()

		// Here is where we convert the settings paths into an ordered list of scopes, and a map of scope paths
		scopes := []string{}
		scopeMap := map[string]string{}
		for _, scope := range b.settings.Paths.Order() {
			path, _ := b.settings.Paths.Get(scope)
			scopes = append(scopes, scope)
			scopeMap[scope] = path
		}
		// all of our files will be .yml files (we will have no prefix)
		suffix := ".yml"

		// So now we can build our file path interpreter based on the captured file path definitions above.
		paths := base_config_provider_file.NewScopedFilePaths(scopes, scopeMap, "", suffix)

		/**
		 * Build a Composite backend based on:
		 *  Connector: a files connecter build using the scope file paths approach
		 *  Usage: All usage (we will change this in the future to include a default scope handler)
		 *  Factory: The config provider factory will be a yaml factory
		 */

		con := base_config_provider_file.NewFileConnector(paths.FilePaths()).Connector() // how to connect to files
		fac := base_config_provider_yaml.NewFactory(con).Factory()                       // how to build Config from the connector
		use := base_config_provider.AllBackendUsage{}.BackendUsage()                     // under what circumstances to use this backend

		backend := base_config_provider.NewCompositeBackend(con, use, fac).Backend()

		tbcp.Add(backend)
	}
	return b.sharedConfigProvider
}

func (b *Builder) ConfigOperations() api.Operations {
	prov := b.configProvider()
	ops := base.NewOperations()

	ops.Add(base_config_provider.NewListOperation(prov).Operation())
	ops.Add(base_config_provider.NewGetOperation(prov).Operation())

	return ops.Operations()
}
