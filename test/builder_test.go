package test

import (
	"os"
	"os/user"
	"path"
	"testing"

	coach_local "github.com/CoachApplication/coach-local"
)

func testBuilder_Builders(t *testing.T) *coach_local.Builder {
	wd, _ := os.Getwd()
	curUser, _ := user.Current()

	localPaths := coach_local.NewSettingScopePaths()
	for _, scope := range []string{"one", "two", "three"} {
		localPaths.Set(scope, path.Join(wd, scope))
	}

	settings := coach_local.Settings{
		ProjectDoesntExist: false,
		User:               *curUser,
		ProjectRootPath:    wd,
		ExecPath:           wd,
		Paths:              *localPaths,
	}

	return coach_local.NewBuilder(settings)
}

func TestBuilder_Operations(t *testing.T) {

}
