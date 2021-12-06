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
	root.GET("/todo/new", actions.NewTask)
	root.POST("/todo", actions.SaveTask)
	root.GET("/todo/{todo_id}/edit", actions.EditTask)
	root.PUT("/todo/{todo_id}", actions.UpdateTask)
	root.DELETE("/todo/{todo_id}", actions.DeleteTask)
	root.ServeFiles("/", base.Assets)
}
