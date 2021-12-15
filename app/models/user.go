package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
	FirstName            string    `json:"first_name" db:"-"`
	LastName             string    `json:"last_name" db:"-"`
	Email                string    `json:"email" db:"email"`
	PasswordHash         string    `json:"password_hash" db:"password_hash"`
	Password             string    `json:"-" db:"-"`
	PasswordConfirmation string    `json:"-" db:"-"`
}

type Users []User

func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(ph)
	return tx.ValidateAndCreate(u)
}

func (u *User) FullName() (string, error) {
	if (strings.TrimSpace(u.FirstName) == "") || (strings.TrimSpace(u.LastName) == "") {
		return "", errors.New("First or last names should not be empty")
	}

	fullName := fmt.Sprintf("%v %v", u.FirstName, u.LastName)
	return fullName, nil
}
