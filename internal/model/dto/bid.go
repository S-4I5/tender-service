package dto

import (
	"github.com/google/uuid"
	"tender-service/internal/model/entity/bid"
	"time"
)

type CreateBidDto struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Status      bid.Status     `json:"status" validate:"required"`
	TenderId    uuid.UUID      `json:"tenderId" validate:"required"`
	AuthorType  bid.AuthorType `json:"authorType" validate:"required"`
	AuthorId    uuid.UUID      `json:"authorId" validate:"required"`
}

type BidDto struct {
	Id          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      bid.Status     `json:"status"`
	TenderId    uuid.UUID      `json:"tenderId"`
	AuthorType  bid.AuthorType `json:"authorType"`
	AuthorId    uuid.UUID      `json:"authorId"`
	Version     int            `json:"version"`
	CreatedAt   time.Time      `json:"createdAt"`
}

type UpdateBidDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
