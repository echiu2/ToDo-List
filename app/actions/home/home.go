package home

import (
	"net/http"
	"time"
	"todolist/app/models"
	"todolist/app/render"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

var (
	// r is a buffalo/render Engine that will be used by actions
	// on this package to render render HTML or any other formats.
	r = render.Engine
)

func Index(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	todos := models.Todos{}
	test := models.Todo{"title", "detail", time.Now(), time.Now()}
	todos = append(todos, test)

	if err := tx.All(&todos); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	c.Set("todos", todos)

	return c.Render(http.StatusOK, r.HTML("home/index.plush.html"))
}
