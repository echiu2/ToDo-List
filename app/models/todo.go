package models

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

type Todo struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Deadline    time.Time `db:"deadline" json:"deadline"`
	IsCompleted bool      `db:"is_completed" json:"is_completed"`
	Title       string    `db:"title" json:"title"`
	Details     string    `db:"details" json:"details"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type Todos []Todo

func (t Todo) TableName() string {
	return "todos"
}

func (t *Todo) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Title, Name: "Title", Message: "Title cannot be empty"},
		&validators.StringLengthInRange{Field: t.Details, Name: "Details", Min: 10, Message: "Details length must be 10 or larger"},
		&validators.TimeIsPresent{Field: t.Deadline, Name: "Deadline", Message: "Deadline not a valid date"},
		&validators.TimeAfterTime{FirstTime: t.Deadline, FirstName: "Deadline", SecondTime: time.Now().Truncate(24 * time.Hour), Message: "Deadline cannot be a past date or current date"},
	)
}
