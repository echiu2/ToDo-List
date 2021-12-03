package actions

import (
	"net/http"
	"todolist/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
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
