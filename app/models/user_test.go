// models/user_test.go
package models_test

import (
	"io/fs"
	"todolist/app/models"

	"github.com/gobuffalo/pop/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Model suite
type ModelSuite struct {
	suite.Suite
	*require.Assertions
	DB       *pop.Connection
	Fixtures fs.FS
}

func (ms *ModelSuite) Test_User_Create() {
	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)

	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	ms.Zero(u.PasswordHash)

	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	ms.NotZero(u.PasswordHash)

	count, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(1, count)
}

func (ms *ModelSuite) Test_User_Create_ValidationErrors() {
	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)

	u := &models.User{
		Password: "password",
	}
	ms.Zero(u.PasswordHash)

	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())

	count, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)
}

func (ms *ModelSuite) Test_User_Create_UserExists() {
	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, count)

	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	ms.Zero(u.PasswordHash)

	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	ms.NotZero(u.PasswordHash)

	count, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(1, count)

	u = &models.User{
		Email:    "mark@example.com",
		Password: "password",
	}
	verrs, err = u.Create(ms.DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())

	count, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(1, count)
}
