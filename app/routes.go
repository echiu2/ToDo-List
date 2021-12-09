package app

import (
	base "todolist"
	"todolist/app/actions"
	"todolist/app/middleware"

	"github.com/gobuffalo/buffalo"
	csrf "github.com/gobuffalo/mw-csrf"
)

// SetRoutes for the application
func setRoutes(root *buffalo.App) {
	root.Use(middleware.Transaction)
	root.Use(middleware.ParameterLogger)
	root.Use(middleware.CSRF)
	root.Use(csrf.New)
	root.Use(actions.Authorize)

	root.GET("/", actions.Index)
	root.Use(actions.SetCurrentUser)
	root.GET("/todo/new", actions.NewTask)
	root.POST("/todo", actions.SaveTask).Name("createTaskPath")
	root.GET("/todo/{todo_id}/edit", actions.EditTask)
	root.PUT("/todo/{todo_id}", actions.UpdateTask)
	root.DELETE("/todo/{todo_id}", actions.DeleteTask)
	root.GET("/users/new", actions.UsersNew)
	root.POST("/users", actions.UsersCreate)
	root.GET("/signin", actions.AuthNew)
	root.POST("/signin", actions.AuthCreate)
	root.DELETE("/signout", actions.AuthDestroy)

	root.Middleware.Skip(actions.Authorize, actions.Index, actions.UsersNew, actions.UsersCreate, actions.AuthNew, actions.AuthCreate, actions.AuthDestroy)
	root.ServeFiles("/", base.Assets)
}
