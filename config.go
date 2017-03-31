package local

import (
	api "github.com/CoachApplication/api"
	base "github.com/CoachApplication/base"
	config "github.com/CoachApplication/config"
	config_provider "github.com/CoachApplication/config/provider"
	config_provider_file "github.com/CoachApplication/config/provider/file"
	config_provider_yaml "github.com/CoachApplication/config/provider/yaml"
)

func (b *Builder) configProvider() config_provider.Provider {
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
		tbcp := config_provider.NewBackendConfigProvider()
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
		paths := config_provider_file.NewScopedFilePaths(scopes, scopeMap, "", suffix)

		/**
		 * Build a Composite backend based on:
		 *  Connector: a files connecter build using the scope file paths approach
		 *  Usage: All usage (we will change this in the future to include a default scope handler)
		 *  Factory: The config provider factory will be a yaml factory
		 */

		con := config_provider_file.NewFileConnector(paths.FilePaths()).Connector() // how to connect to files
		fac := config_provider_yaml.NewFactory(con).Factory()                       // how to build Config from the connector
		use := (&config_provider.AllBackendUsage{}).BackendUsage()                  // under what circumstances to use this backend

		backend := config_provider.NewCompositeBackend(con, use, fac).Backend()

		tbcp.Add(backend)
	}
	return b.sharedConfigProvider
}

// ConfigWrapper build a config.Wrapper from the builder ConfigOperations
func (b *Builder) configWrapper() config.Wrapper {
	return config.NewStandardWrapper(b.configOperations(), b.context).Wrapper()
}

func (b *Builder) configOperations() api.Operations {
	prov := b.configProvider()
	ops := base.NewOperations()

	ops.Add(config_provider.NewListOperation(prov).Operation())
	ops.Add(config_provider.NewGetOperation(prov).Operation())

	return ops.Operations()
}
