package local

import (
	coach_api "github.com/james-nesbitt/coach-api"
	coach_base "github.com/james-nesbitt/coach-base"
)

func (b *Builder) ConfigOperations() coach_api.Operations {
	return coach_base.NewOperations()
}


type ConfigResolver struct {

}

func (cr *ConfigResolver) ListScopes