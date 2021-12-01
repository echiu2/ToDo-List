package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func Index(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("/todos/index.plush.html"))
}
