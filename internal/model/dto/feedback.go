package dto

import (
	"github.com/google/uuid"
	"time"
)

type FeedbackDto struct {
	Id          uuid.UUID
	Description string
	CreatedAt   time.Time
}
