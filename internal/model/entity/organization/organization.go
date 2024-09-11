package organization

import (
	"github.com/google/uuid"
	"time"
)

type Type string

const (
	IE  Type = "IE"
	LLC Type = "LLC"
	JSC Type = "JSC"
)

type Organization struct {
	Id          uuid.UUID
	Name        string
	Description string
	Type        Type
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
