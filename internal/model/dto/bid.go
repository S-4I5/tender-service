package dto

import (
	"github.com/google/uuid"
	"tender-service/internal/model/entity/bid"
	"time"
)

type CreateBidDto struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Status      bid.Status     `yaml:"status"`
	TenderId    uuid.UUID      `yaml:"tenderId"`
	AuthorType  bid.AuthorType `json:"authorType"`
	//OrganizationId uuid.UUID  `yaml:"organizationId"`
	AuthorId uuid.UUID `yaml:"authorId"`
}

type BidDto struct {
	Id          uuid.UUID      `yaml:"id"`
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Status      bid.Status     `yaml:"status"`
	TenderId    uuid.UUID      `yaml:"tenderId"`
	AuthorType  bid.AuthorType `yaml:"authorType"`
	AuthorId    uuid.UUID      `yaml:"authorId"`
	Version     int            `yaml:"version"`
	CreatedAt   time.Time      `yaml:"createdAt"`
}

type UpdateBidDto struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}
