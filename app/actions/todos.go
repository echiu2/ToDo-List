package actions

import (
	"fmt"
	"net/http"
	"time"
	"todolist/app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func Index(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	todos := models.Todos{}
	uid := c.Session().Get("current_user_id")

	if err := tx.Where("user_id = (?)", uid).All(&todos); err != nil {
		return c.Error(http.StatusNotFound, errors.Wrap(err, "Index - error while finding todo object"))
	}

	c.Set("todos", todos)

	return c.Render(http.StatusOK, r.HTML("/todos/index.plush.html"))
}

func NewTask(c buffalo.Context) error {
	time := time.Now()
	todo := models.Todo{Deadline: time}

	c.Set("todo", todo)

	return c.Render(http.StatusOK, r.HTML("todos/new.plush.html"))
}

func SaveTask(c buffalo.Context) error {
	todo := models.Todo{}

	if err := c.Bind(&todo); err != nil {
		return c.Error(http.StatusInternalServerError, errors.Wrap(err, "Save - error while binding to Todo"))
	}

	if verrs := todo.Validate(); verrs.HasAny() {
		c.Set("todo", todo)
		c.Set("errors", verrs.Errors)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("todos/new.plush.html"))
	}

	uid := c.Session().Get("current_user_id")
	userID := fmt.Sprintf("%v", uid)

	todo.UserID = uuid.FromStringOrNil(userID)

	tx := c.Value("tx").(*pop.Connection)

	err := tx.Create(&todo)
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.Wrap(err, "Save - error creating a new todo object"))
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func EditTask(c buffalo.Context) error {
	todoID := c.Param("todo_id")
	todo := models.Todo{}

	tx := c.Value("tx").(*pop.Connection)

	if err := tx.Find(&todo, todoID); err != nil {
		return c.Error(http.StatusNotFound, errors.Wrap(err, "Edit - error while finding todo object"))
	}

	c.Set("todo", todo)

	return c.Render(http.StatusOK, r.HTML("todos/edit.plush.html"))
}

func UpdateTask(c buffalo.Context) error {
	todoID := c.Param("todo_id")
	todo := models.Todo{}

	tx := c.Value("tx").(*pop.Connection)

	if err := tx.Find(&todo, todoID); err != nil {
		return c.Error(http.StatusNotFound, errors.Wrap(err, "Edit - error while finding todo object"))
	}

	if err := c.Bind(&todo); err != nil {
		return c.Error(http.StatusInternalServerError, errors.Wrap(err, "Edit - error while binding to Todo"))
	}

	if verrs := todo.Validate(); verrs.HasAny() {
		c.Set("todo", todo)
		c.Set("errors", verrs.Errors)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("todos/edit.plush.html"))
	}

	err := tx.Update(&todo)
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.Wrap(err, "Update - error updating an existing todo object"))
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func DeleteTask(c buffalo.Context) error {
	todoID := c.Param("todo_id")
	todo := models.Todo{}

	tx := c.Value("tx").(*pop.Connection)

	if err := tx.Find(&todo, todoID); err != nil {
		return c.Error(http.StatusNotFound, errors.Wrap(err, "Delete - error while finding todo object"))
	}

	err := tx.Destroy(&todo)
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.Wrap(err, "Delete - error while trying to delete an existing todo object"))
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
