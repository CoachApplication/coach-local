package local

import (
	api "github.com/CoachApplication/api"
	handler_dockercli_configwrapper "github.com/CoachApplication/handler-dockercli/configwrapper"
)

/**
 * Orchestration is handled via Docker stacks, using the dockercli handler
 */

func (b *Builder) orchestrateOperations() api.Operations {
	wr := b.configWrapper()
	return handler_dockercli_configwrapper.MakeOrchestrateOperations(wr)
}
