package entity

import (
	"github.com/google/uuid"
	"time"
)

type Employee struct {
	Id        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
