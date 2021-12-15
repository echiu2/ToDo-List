package actions_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"todolist/app/models"

	"github.com/gofrs/uuid"
)

func (as ActionSuite) Test_Index_Status_OK() {
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	as.Session.Set("current_user_id", userA.ID)
	res := as.HTML("/").Get()
	as.Equal(http.StatusOK, res.Code)

	as.Session.Clear()
	res = as.HTML("/").Get()
	as.Equal(http.StatusOK, res.Code)

	body := res.Body.String()

	as.Contains(body, "sign in")
	as.Contains(body, "register")
}

func (as ActionSuite) Test_Navbar_Text() {
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	as.Session.Set("current_user_id", userA.ID)
	res := as.HTML("/").Get()

	body := res.Body.String()

	time := time.Now()
	message := fmt.Sprintf("%v", time.Format("Monday 02, January 2006"))

	as.Contains(body, "ToDo List")
	as.Contains(body, "Incomplete Tasks")
	as.Contains(body, "Completed Tasks")
	as.Contains(body, "Sign out")
	as.Contains(body, message)
}

func (as ActionSuite) Test_Welcome_Text() {
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	userB := models.User{Email: "jarias@testing.com"}
	as.NoError(as.DB.Create(&userB)) // we need the user to be created in the DB

	as.Session.Set("current_user_id", userA.ID)
	res := as.HTML("/").Get()

	body := res.Body.String()

	as.Contains(body, "Add Task")

	as.Contains(body, "Welcome, "+userA.Email)
	as.NotContains(body, "Welcome, "+userB.Email)
}

func (as ActionSuite) Test_Empty_Todo_List() {
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	todos := models.Todos{}

	as.Session.Set("current_user_id", userA.ID)
	res := as.HTML("/").Get()

	body := res.Body.String()

	as.Empty(todos)

	as.NotContains(body, `<th scope="col">Task</th>`)
	as.NotContains(body, `<th scope="col">Complete By</th>`)
	as.NotContains(body, `<th scope="col">Actions</th>`)
	as.Contains(body, "No current tasks needs to be done at the moment.")
}

func (as ActionSuite) Test_Populated_Todo_List() {
	// 1. Arrange
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	todos := models.Todos{
		{
			Title:       "Test Title 1",
			IsCompleted: false,
			Details:     "This is a long chain of details",
			UserID:      userA.ID,
			Deadline:    time.Date(2020, 12, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:       "Test Title 2",
			IsCompleted: false,
			Details:     "This is another long chain of details",
			UserID:      userA.ID,
			Deadline:    time.Date(2020, 12, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:       "Test Title 3",
			IsCompleted: true,
			Details:     "This is yet another long chain of details",
			UserID:      userA.ID,
			Deadline:    time.Date(2020, 12, 8, 0, 0, 0, 0, time.UTC),
		},
	}

	as.NoError(as.DB.Create(&todos))
	as.Session.Set("current_user_id", userA.ID) // this is to simulate the user is logged into the system

	// 2. Act
	res := as.HTML("/").Get()
	as.Equal(http.StatusOK, res.Code) // we are saying we should obtain a 200 from "/""

	// 3. Assert
	body := res.Body.String()

	as.NotEmpty(todos)

	for _, v := range todos {
		as.Contains(body, v.Title)
		as.Contains(body, v.Deadline.Format("02 Jan 2006"))
	}

	as.Contains(body, `<th scope="col">Task</th>`)
	as.Contains(body, `<th scope="col">Complete By</th>`)
	as.Contains(body, `<th scope="col">Actions</th>`)
}

func (as ActionSuite) Test_NewTask_Status_OK() {
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	as.Session.Set("current_user_id", userA.ID)
	res := as.HTML("/todo/new").Get()
	as.Equal(http.StatusOK, res.Code)

	as.Session.Clear()
	res = as.HTML("/").Get()
	as.Equal(http.StatusOK, res.Code)

	body := res.Body.String()

	as.Contains(body, "sign in")
	as.Contains(body, "register")
}

func (as ActionSuite) Test_NewTask_Initial_Form() {
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	as.Session.Set("current_user_id", userA.ID)

	res := as.HTML("/todo/new").Get()
	as.Equal(http.StatusOK, res.Code)

	body := res.Body.String()
	time := time.Now().Format("2006-01-02")

	as.Contains(body, "New Task")
	as.Contains(body, "Title")
	as.Contains(body, "Limit Date")
	as.Contains(body, "Details")
	as.Contains(body, time)

	as.Contains(body, "Cancel")
	as.Contains(body, "Create Task")
}

func (as ActionSuite) Test_Submit_NewTask_Form() {
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	as.Session.Set("current_user_id", userA.ID)

	task := models.Todo{
		Deadline:    time.Date(2021, 12, 22, 0, 0, 0, 0, time.UTC),
		IsCompleted: false,
		Title:       "Test123",
		Details:     "Test123456",
		UserID:      userA.ID,
	}

	res := as.HTML("/todo/").Post(&task)
	as.Equal(http.StatusSeeOther, res.Code)
	as.NoError(as.DB.Create(&task))
	as.NoError(as.DB.Where("id = ?", task.ID).First(&models.Todo{}))

	notValidTasks := []url.Values{
		{
			"Deadline":    []string{"2021-12-14"},
			"IsCompleted": []string{"false"},
			"Title":       []string{"test123"},
			"Details":     []string{"test1"},
		},
		{
			"Deadline":    []string{"2021-12-14"},
			"IsCompleted": []string{"false"},
			"Title":       []string{""},
			"Details":     []string{"test123456"},
		},
	}

	for _, v := range notValidTasks {
		check := as.HTML("/todo").Post(v)
		as.Equal(http.StatusUnprocessableEntity, check.Code)
	}

}

func (as ActionSuite) Test_EditTask_Initial_Form() {
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	userB := models.User{Email: "test@testing.com"}
	as.NoError(as.DB.Create(&userB))

	as.Session.Set("current_user_id", userA.ID)

	taskA := models.Todo{
		Deadline:    time.Now(),
		IsCompleted: false,
		Title:       "Test123",
		Details:     "Test123456",
		UserID:      userA.ID,
	}
	as.NoError(as.DB.Create(&taskA))

	taskB := models.Todo{
		Deadline:    time.Now(),
		IsCompleted: false,
		Title:       "Test123",
		Details:     "Test123456",
		UserID:      userB.ID,
	}
	as.NoError(as.DB.Create(&taskB))

	res := as.HTML("/todo/%v/edit", taskA.ID).Get()
	as.Equal(http.StatusOK, res.Code)

	body := res.Body.String()
	time := time.Now().Format("2006-01-02")

	as.Contains(body, "Edit Task")
	as.Contains(body, "Title")
	as.Contains(body, "Limit Date")
	as.Contains(body, "Details")
	as.Contains(body, time)

	as.Contains(body, "Cancel")
	as.Contains(body, "Edit Task")

	as.Contains(body, taskA.Title)
	as.Contains(body, taskA.Deadline.Format("Monday 02, January 2006"))
	as.Contains(body, taskA.Details)

	res = as.HTML("/todo/%v/edit", taskB.ID).Get()
	as.Equal(http.StatusForbidden, res.Code)
	as.Error(errors.New("GuardEditMW - you do not have access and authorization to this url or task"))

}

func (as ActionSuite) Test_Submit_Update_Form() {
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	as.Session.Set("current_user_id", userA.ID)

	taskA := models.Todo{
		Deadline:    time.Now(),
		IsCompleted: false,
		Title:       "Test123",
		Details:     "Test123456",
		UserID:      userA.ID,
	}
	as.NoError(as.DB.Create(&taskA))

	form := url.Values{
		"Deadline":    []string{taskA.Deadline.String()},
		"IsCompleted": []string{strconv.FormatBool(taskA.IsCompleted)},
		"Title":       []string{taskA.Title},
		"Details":     []string{taskA.Details},
	}

	res := as.HTML("/todo/%v/", taskA.ID).Get()
	as.Equal(http.StatusNotFound, res.Code)

	res = as.HTML("/todo/%v/", taskA.ID).Put(form)
	as.Equal(http.StatusSeeOther, res.Code)
}

func (as ActionSuite) Test_Delete_Task() {
	userA := models.User{Email: "echiu@testing.com"}
	as.NoError(as.DB.Create(&userA)) // we need the user to be created in the DB

	as.Session.Set("current_user_id", userA.ID)

	taskA := models.Todo{
		ID:          uuid.Must(uuid.NewV4()),
		Deadline:    time.Now(),
		IsCompleted: false,
		Title:       "Test123",
		Details:     "Test123456",
		UserID:      userA.ID,
	}
	as.NoError(as.DB.Create(&taskA))

	res := as.HTML("/todo/%v/", taskA.ID).Delete()
	as.Equal(http.StatusSeeOther, res.Code)

	count, _ := as.DB.Count(&taskA)
	as.Equal(count, 0)
}
