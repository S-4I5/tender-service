package dto

import (
	"github.com/google/uuid"
	"time"
)

type FeedbackDto struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
