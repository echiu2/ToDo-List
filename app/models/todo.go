package models

import (
	"time"
)

type Todo struct {
	Title     string    `db:"title" json:"title"`
	Details   string    `db:"details" json:"details"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Todos []Todo

func (t Todo) TableName() string {
	return "todos"
}
