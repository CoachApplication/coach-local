package local_test

import (
	"testing"

	"context"
	config "github.com/CoachApplication/config"
	local "github.com/CoachApplication/handler-local"
	"time"
)

func TestBuilder_Operations(t *testing.T) {
	b := local.MakeLocalBuilder(t, nil)

	b.Activate([]string{"config"}, nil)

	ops := b.Operations()
	list := ops.Order()

	if len(list) == 0 {
		t.Error("LocalBuilder returned no operations")
	} else if !isInSlice(config.OPERATION_ID_CONFIG_GET, list) {
		t.Error("LocalBuilder did not provide a config.get operation")
	} else if !isInSlice(config.OPERATION_ID_CONFIG_LIST, list) {
		t.Error("LocalBuilder did not provide a config.list operation")
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
	} else if scA, err := cw.Get("A"); err != nil {
		t.Error("ConfigWrapper produced an error during Get()", err, ls)
	} else {
		var ts local.TestStruct

		if dcA, err := scA.Get("one"); err != nil {
			t.Error("ConfigWrapper ScopedConfigList returned an error when retrieving a valid Config: ", err.Error())
		} else {
			res := dcA.Get(&ts)
			<-res.Finished()

			if !res.Success() {
				t.Error("ConfigWrapper ScopedConfigList returned an failure when unmarshalling a Config: ", res.Errors())
			}
		}
	}

	if ls, err := cw.List(); err != nil {
		t.Error("ConfigWrapper from LocalBuilder Operations() gave an error on List()", err)
	} else if len(ls) == 0 {
		t.Error("ConfigWrapper from LocalBuilder Returned no Config keys", ls)
	} else if scA, err := cw.Get("A"); err != nil {
		t.Error("ConfigWrapper produced an error during Get()", err, ls)
	} else {
		scopes := scA.Order()
		t.Log("Discovered scopes :", scopes)

		var ts local.TestStruct

		if dcA, err := scA.Get(config.CONFIG_SCOPE_DEFAULT); err != nil {
			t.Error("ConfigWrapper ScopedConfigList returned an error when retrieving the default Config: ", err.Error())
		} else {
			res := dcA.Get(&ts)
			<-res.Finished()

			if !res.Success() {
				t.Error("ConfigWrapper ScopedConfigList returned an failure when unmarshalling the default Config: ", res.Errors())
			}
		}
	}
}

func isInSlice(vl string, sl []string) bool {
	for _, each := range sl {
		if each == vl {
			return true
		}
	}
	return false
}
