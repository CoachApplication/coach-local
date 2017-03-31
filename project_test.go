package local_test

import (
	"testing"

	local "github.com/CoachApplication/handler-local"
	project "github.com/CoachApplication/project"
)

func TestBuilder_Operations_Project(t *testing.T) {
	b := local.MakeLocalBuilder(t, nil)

	b.Activate([]string{"project"}, nil)

	ops := b.Operations()
	list := ops.Order()

	if len(list) == 0 {
		t.Error("LocalBuilder:project returned no operations")
	} else {
		for _, id := range []string{
			project.OPERATION_ID_NAME,
			project.OPERATION_ID_LABEL,
			project.OPERATION_ID_ENV_MAP,
			project.OPERATION_ID_ENV_MAP,
			project.OPERATION_ID_LABEL,
			project.OPERATION_ID_LABEL,
		} {
			if !isInSlice(id, list) {
				t.Errorf("LocalBuilder:project did not provide a %s operation", id)
			}
		}
	}
}


