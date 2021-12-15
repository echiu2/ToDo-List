package actions_test

import (
	"testing"
	"todolist/app"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/suite/v3"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	bapp := app.New()

	pop.Debug = false
	as := &ActionSuite{suite.NewAction(bapp)}
	suite.Run(t, as)
}
