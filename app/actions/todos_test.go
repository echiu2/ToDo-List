package actions

import (
	"todolist/app/models"

	"github.com/gofrs/uuid"
)

func (as ActionSuite) Test_Todos_Index() {
	as.Session.Set("current_user_id", "7e2d187a-de99-4e8b-817b-409edee91141") // this is to simulate the user is logged into the system

	user := models.User{
		ID: uuid.FromStringOrNil("7e2d187a-de99-4e8b-817b-409edee91141")
		Email: "echiu@testing.com",
	}

	as.NoError(as.DB.Create(&user)) // we need the user to be created in the DB

	res := as.HTML("/").Get()
	as.Equal(http.StatusOK, res.Code) // we are saying we should obtain a 200 from "/""

	//let's run this thing
	// ox test ./... -count=1
	// 
}

// yup, theres some noise. I will write then explain : )
// 