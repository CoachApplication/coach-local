package local

import (
	api "github.com/CoachApplication/api"
	project_configwrapper "github.com/CoachApplication/project/configwrapper"
)

func (b *Builder) projectOperations() api.Operations {
	wr := b.configWrapper()
	return project_configwrapper.MakeOperations(wr)
}
