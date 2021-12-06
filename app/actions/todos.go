package actions

import (
	"net/http"
	"todolist/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func Index(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	todos := models.Todos{}

	if err := tx.All(&todos); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	c.Set("todos", todos)

	return c.Render(http.StatusOK, r.HTML("/todos/index.plush.html"))
}

func NewTask(c buffalo.Context) error {
	todo := models.Todo{}
	c.Set("todo", todo)

	return c.Render(http.StatusOK, r.HTML("todos/new.plush.html"))
}

func PostNewTask(c buffalo.Context) error {
	u, _ := uuid.NewV4()
	todo := models.Todo{ID: u}

	if err := c.Bind(&todo); err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)

	err := tx.Create(&todo)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("todo", todo)

	return c.Render(http.StatusOK, r.HTML("todos/new.plush.html"))
}
