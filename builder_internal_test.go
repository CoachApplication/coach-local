package local

import (
	"context"
	"os"
	"os/user"
	"path"
	"testing"

	config_provider "github.com/CoachApplication/config/provider"
	config_provider_file "github.com/CoachApplication/config/provider/file"
	config_provider_yaml "github.com/CoachApplication/config/provider/yaml"
)

// Centralize making of test local.Settings so that we can be sure of them
func MakeLocalSettings(t *testing.T) Settings {
	wd, _ := os.Getwd()
	curUser, _ := user.Current()

	t.Logf("Builder Base [user: %s][wd: %s]", curUser.Name, wd)

	localPaths := NewSettingScopePaths()
	for _, scope := range []string{"default", "user", "project"} {
		scopePath := path.Join(wd, "test", "config", scope)
		localPaths.Set(scope, scopePath)
		t.Logf("Builder Base Added scope: [scope: %s][path: %s]", scope, scopePath)
	}

	return Settings{
		ProjectDoesntExist: false,
		User:               *curUser,
		ProjectRootPath:    wd,
		ExecPath:           wd,
		Paths:              *localPaths,
	}
}

func MakeLocalBuilder(t *testing.T, ctx context.Context) *Builder {
	settings := MakeLocalSettings(t)
	if ctx == nil {
		ctx = context.Background()
	}
	return NewBuilder(ctx, settings)
}

type TestStruct struct {
	Integers []int `yaml:"integers"`
}

// Here we test the process that the builder uses to build a config.Provider
func TestBuilder_configProviderSteps(t *testing.T) {
	settings := MakeLocalSettings(t)

	// This provider
	prov := config_provider.NewBackendConfigProvider()

	// Here is where we convert the settings paths into an ordered list of scopes, and a map of scope paths
	scopes := []string{}
	scopeMap := map[string]string{}
	for _, scope := range settings.Paths.Order() {
		p, _ := settings.Paths.Get(scope)
		scopes = append(scopes, scope)
		scopeMap[scope] = p
		t.Log("Added scope path: ", scope, p)
	}
	suffix := ".yml"

	paths := config_provider_file.NewScopedFilePaths(scopes, scopeMap, "", suffix)

	con := config_provider_file.NewFileConnector(paths.FilePaths()).Connector() // how to connect to files
	fac := config_provider_yaml.NewFactory(con).Factory()                       // how to build Config from the connector
	use := (&config_provider.AllBackendUsage{}).BackendUsage()                  // under what circumstances to use this backend

	backend := config_provider.NewCompositeBackend(con, use, fac).Backend()
	prov.Add(backend)

	t.Log("Paths:Keys()", paths.Keys())
	t.Log("Paths:Scopes()", paths.Scopes())
	t.Log("Backend:Keys()", con.Keys())
	t.Log("Backend:Scopes()", con.Scopes())
	t.Log("Backend:Keys()", backend.Keys())
	t.Log("Backend:Scopes()", backend.Scopes())
	t.Log("Provider:Keys()", prov.Keys())
	t.Log("Provider:Scopes()", prov.Scopes())

	if len(prov.Keys()) != 3 {
		t.Error("TestBuilder filepaths returned the wrong number of keys: ", prov.Keys())
	}
	if len(prov.Scopes()) != 3 {
		t.Error("TestBuilder filepaths returned the wrong number of scopes: ", prov.Scopes())
	}
}

// Here we test the process that the builder uses to build a config.Provider
func TestBuilder_configProvider(t *testing.T) {
	b := MakeLocalBuilder(t, nil)
	prov := b.configProvider()

	t.Log("Provider:Keys()", prov.Keys())
	t.Log("Provider:Scopes()", prov.Scopes())

	if len(prov.Keys()) != 3 {
		t.Error("TestBuilder config.Provider returned the wrong number of keys:", prov.Keys())
	}
	if len(prov.Scopes()) != 3 {
		t.Error("TestBuilder config.Provider returned the wrong number of keys:", prov.Keys())
	}

	var valB TestStruct
	if scB, err := prov.Get("integers", "default"); err != nil {
		t.Error("Testbuilder config.Provider returned an error retrieving valid key scope pair")
	} else {
		res := scB.Get(&valB)
		<-res.Finished()

		if !res.Success() {
			t.Error("TestBuilder config.Provider gave an error marshalling config to TestStruct:", res.Errors())
		} else if &valB == nil {
			t.Error("TestBuilder config.Provider failed to marshal config properly.  It is still nil")
		} else if valAInts := valB.Integers; len(valAInts) == 0 {
			t.Error("TestBuilder config.Provider failed to marshal config properly.  Not enough integers")
		} else if valAInts[0] != 1 {
			t.Error("TestBuilder config.Provider failed to marshal config properly.  Wrong integers value")
		}
	}
}
