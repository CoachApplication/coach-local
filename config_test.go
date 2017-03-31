package local_test

import (
	"context"
	"testing"
	"time"

	config "github.com/CoachApplication/config"
	local "github.com/CoachApplication/handler-local"
)

func TestBuilder_Operations_Config(t *testing.T) {
	b := local.MakeLocalBuilder(t, nil)

	b.Activate([]string{"config"}, nil)

	ops := b.Operations()
	list := ops.Order()

	if len(list) == 0 {
		t.Error("LocalBuilder:config returned no operations")
	} else {

		for _, id := range []string{
			config.OPERATION_ID_GET,
			config.OPERATION_ID_LIST,
		} {
			if !isInSlice(id, list) {
				t.Errorf("LocalBuilder:config did not provide a %s operation", id)
			}
		}
	}
}

func TestBuilder_ConfigOperations(t *testing.T) {
	b := local.MakeLocalBuilder(t, nil)
	b.Activate([]string{"config"}, nil)

	ops := b.Operations()

	if ls := ops.Order(); len(ls) == 0 {
		t.Error("LocalBuilder Returned no Operation keys", ls)
	}
}

func TestBuilder_ConfigWrapper(t *testing.T) {
	b := local.MakeLocalBuilder(t, nil)
	b.Activate([]string{"config"}, nil)

	dur, _ := time.ParseDuration("2s")
	ctx, _ := context.WithTimeout(context.Background(), dur)
	cw := config.NewStandardWrapper(b.Operations(), ctx)

	if ls, err := cw.List(); err != nil {
		t.Error("ConfigWrapper from LocalBuilder Operations() gave an error on List()", err)
	} else if len(ls) == 0 {
		t.Error("ConfigWrapper from LocalBuilder Returned no Config keys", ls)
	} else if scB, err := cw.Get("integers"); err != nil {
		t.Error("ConfigWrapper produced an error during Get()", err, ls)
	} else {
		var ts local.TestStruct

		if dcB, err := scB.Get("default"); err != nil {
			t.Error("ConfigWrapper ScopedConfigList returned an error when retrieving a valid Config: ", err.Error())
		} else {
			res := dcB.Get(&ts)
			<-res.Finished()

			if !res.Success() {
				t.Error("ConfigWrapper ScopedConfigList returned an failure when unmarshalling a Config: ", res.Errors())
			}
		}
	}
}