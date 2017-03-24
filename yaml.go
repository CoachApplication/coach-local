package local

import (
	"bytes"
	yaml "gopkg.in/yaml.v2"
	"io"

	api "github.com/james-nesbitt/coach-api"
	base "github.com/james-nesbitt/coach-base"
	base_configprovider "github.com/james-nesbitt/coach-base/handler/configprovider"
	base_config "github.com/james-nesbitt/coach-base/operation/configuration"
)

/**
 * Here we have the configprovider related architecture that will convert yaml files to config
 */

// YamlBackend a configprofider backend that can make Config objects from a connector
type YamlBackend struct {
	connector base_configprovider.Connector
	usage     base_configprovider.BackendUsage
}

// NewYamlBackend Constructor for YamlBackend
func NewYamlBackend(con base_configprovider.Connector, us base_configprovider.BackendUsage) *YamlBackend {
	return &YamlBackend{
		connector: con,
		usage:     us,
	}
}

func (yb *YamlBackend) Handles(key, scope string) bool {
	return yb.usage.Handles(key, scope)
}
func (yb *YamlBackend) Scopes() []string {
	return yb.connector.Scopes()
}
func (yb *YamlBackend) Keys() []string {
	return yb.connector.Keys()
}
func (yb *YamlBackend) Get(key, scope string) (base_config.Config, error) {

	/** @TODO should we check if the path exists first?  If we do then you can't write to new configs? */

	return NewYamlConnectorConfig(key, scope, yb.connector).Config(), nil
}

// YamlConfig Build Config by marshalling yml from a connector
type YamlConnectorConfig struct {
	key       string
	scope     string
	connector base_configprovider.Connector
}

func NewYamlConnectorConfig(key, scope string, con base_configprovider.Connector) *YamlConnectorConfig {
	return &YamlConnectorConfig{
		key:       key,
		scope:     scope,
		connector: con,
	}
}

// Marshall gets a configuration and apply it to a target struct
func (ycc *YamlConnectorConfig) Config() base_config.Config {
	return base_config.Config(ycc)
}

// Marshall gets a configuration and apply it to a target struct
func (ycc *YamlConnectorConfig) Get(target interface{}) api.Result {

	res := base.NewResult()

	go func(key, scope string, con base_configprovider.Connector) {
		defer res.MarkFinished()

		if r, err := con.Get(key, scope); err != nil {
			res.AddError(err)
			res.MarkFailed()
		} else {
			defer r.Close()
			buf := bytes.Buffer{}
			if _, err := buf.ReadFrom(r); err != nil {
				res.AddError(err)
				res.MarkFailed()
			} else if err := yaml.Unmarshal(buf.Bytes(), target); err != nil {
				res.AddError(err)
				res.MarkFailed()
			} else {
				res.MarkSucceeded()
			}
		}
	}(ycc.key, ycc.scope, ycc.connector)

	return res.Result()
}

// UnMarshall sets a Config value by converting a passed struct into a configuration
// The expects that the values assigned are permanently saved
func (ycc *YamlConnectorConfig) Set(source interface{}) api.Result {
	res := base.NewResult()

	go func(key, scope string, con base_configprovider.Connector) {
		defer res.MarkFinished()

		if b, err := yaml.Marshal(source); err != nil {
			res.AddError(err)
			res.MarkFailed()
		} else {
			rc := io.ReadCloser(&readCloserWrapper{Reader: bytes.NewBuffer(b)})
			if err := con.Set(key, scope, rc); err != nil {
				res.AddError(err)
				res.MarkFailed()
			} else {
				res.MarkSucceeded()
			}
		}
	}(ycc.key, ycc.scope, ycc.connector)

	return res.Result()
}

// Just a simple io.Reader wrapper to make it also a ReadCloser
type readCloserWrapper struct {
	io.Reader
}

// Close actually does nothing, so make sure that your reader doesn't need closing
func (rcw *readCloserWrapper) Close() error {
	return nil
}
