package app

import (
	base "todolist"
	"todolist/app/actions"
	"todolist/app/middleware"

	"github.com/gobuffalo/buffalo"
)

// SetRoutes for the application
func setRoutes(root *buffalo.App) {
	root.Use(middleware.Transaction)
	root.Use(middleware.ParameterLogger)
	root.Use(middleware.CSRF)

	root.GET("/", actions.Index)
	root.ServeFiles("/", base.Assets)
}
