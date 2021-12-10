package app

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
)

var (
	root *buffalo.App
	ENV  = envy.Get("GO_ENV", "development")
)

// App creates a new application with default settings and reading
// GO_ENV. It calls setRoutes to setup the routes for the app that's being
// created before returning it
func New() *buffalo.App {
	if root != nil {
		return root
	}

	pop.Debug = true

	root = buffalo.New(buffalo.Options{
		Env:         ENV,
		SessionName: "_todolist_session",
	})

	// Setting the routes for the app
	setRoutes(root)

	return root
}
