// middleware package is intended to host the middlewares used
// across the app.
package middleware

import (
	"fmt"
	"net/http"
	"todolist/app/models"

	"github.com/gobuffalo/buffalo"
	tx "github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	csrf "github.com/gobuffalo/mw-csrf"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

var (
	// Transaction middleware wraps the request with a pop
	// transaction that is committed on success and rolled
	// back when errors happen.
	Transaction = tx.Transaction(models.DB())

	// ParameterLogger logs out parameters that the app received
	// taking care of sensitive data.
	ParameterLogger = paramlogger.ParameterLogger

	// CSRF middleware protects from CSRF attacks.
	CSRF = csrf.New
)

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		fmt.Println("SetCurrenUser MW")
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		fmt.Println("Authorize MW")
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}

func IncompletedTaskMW(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		fmt.Println("IncompletedTask MW")
		incompletedTask := 0
		if uid := c.Session().Get("current_user_id"); uid != nil {
			todos := models.Todos{}
			tx := c.Value("tx").(*pop.Connection)
			if err := tx.Where("user_id = (?) AND is_completed = false", uid).All(&todos); err != nil {
				// note -- correct error management
				return c.Error(http.StatusNotFound, errors.Wrap(err, "IncompletedTaskMW - error while finding user's todo object"))
			}
			incompletedTask = len(todos)
			c.Set("incompletedTask", incompletedTask)
		}

		return next(c)
	}
}

func GuardEdit(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		fmt.Println("GuardEditMW")

		uid := c.Session().Get("current_user_id")
		if uid == nil {
			return c.Error(http.StatusForbidden, errors.New("GuardEditMW - current_user_id is not set"))
		}

		todo := models.Todo{}
		todoID := c.Param("todo_id")
		tx := c.Value("tx").(*pop.Connection)

		if err := tx.Find(&todo, todoID); err != nil {
			return c.Error(http.StatusNotFound, errors.Wrap(err, "GuardEditMW - error while finding user's todo object"))
		}

		if todo.UserID != uid {
			return c.Error(http.StatusForbidden, errors.New("GuardEditMW - you do not have access and authorization to this url or task"))
		}

		return next(c)
	}
}
