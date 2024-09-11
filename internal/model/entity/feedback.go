package entity

import (
	"github.com/google/uuid"
	"time"
)

type Feedback struct {
	Id          uuid.UUID
	BidId       uuid.UUID
	Description string
	Username    string
	CreatedAt   time.Time
}
