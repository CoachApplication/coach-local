package local_test

import (
	"testing"
	local "github.com/CoachApplication/handler-local"
	handler_dockercli_stack "github.com/CoachApplication/handler-dockercli/stack"
)

func TestBuilder_Operations_Dockercli(t *testing.T) {
	b := local.MakeLocalBuilder(t, nil)

	b.Activate([]string{"dockercli"}, nil)

	ops := b.Operations()
	list := ops.Order()

	if len(list) == 0 {
		t.Error("LocalBuilder:dockercli returned no operations")
	} else {
		for _, id := range []string{
			handler_dockercli_stack.OPERATION_ID_ORCHESTRATE_UP,
			handler_dockercli_stack.OPERATION_ID_ORCHESTRATE_DOWN,
		} {
			if !isInSlice(id, list) {
				t.Errorf("LocalBuilder:dockercli did not provide a %s operation", id)
			}
		}
	}
}
git add ./