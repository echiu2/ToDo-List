package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Todo struct {
	ID          uuid.UUID `db:"id" json:"id"`
	LimitDate   time.Time `db:"limit_date" json:"limit_date"`
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
