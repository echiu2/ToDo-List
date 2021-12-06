package actions

import (
	"net/http"
	"time"
	"todolist/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

func Index(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	todos := models.Todos{}

	if err := tx.All(&todos); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	time := time.Now()

	c.Set("todos", todos)
	c.Set("time", time)

	return c.Render(http.StatusOK, r.HTML("/todos/index.plush.html"))
}

func NewTask(c buffalo.Context) error {
	todo := models.Todo{}
	time := time.Now()

	c.Set("todo", todo)
	c.Set("time", time)

	return c.Render(http.StatusOK, r.HTML("todos/new.plush.html"))
}

func Save(c buffalo.Context) error {
	todo := models.Todo{}

	if err := c.Bind(&todo); err != nil {
		return c.Error(http.StatusInternalServerError, errors.Wrap(err, "Save - error while binding to Todo"))
	}

	tx := c.Value("tx").(*pop.Connection)

	err := tx.Create(&todo)
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.Wrap(err, "Save - error creating a new todo object"))
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
