package bid

import (
	"github.com/google/uuid"
	"time"
)

type AuthorType string

const (
	AuthorOrganization AuthorType = "Organization"
	AuthorUser         AuthorType = "User"
)

type Status string

const (
	Created   Status = "Created"
	Published Status = "Published"
	Canceled  Status = "Canceled"
	Approved  Status = "Approved"
	Rejected  Status = "Rejected"
)

func (b *Bid) IsVisible() bool {
	return b.Status == Published || b.Status == Approved || b.Status == Rejected
}

func IsSelectableByOwner(status Status) bool {
	return status == Created || status == Published || status == Canceled
}

type Bid struct {
	Id          uuid.UUID
	Name        string
	Description string
	Status      Status
	TenderId    uuid.UUID
	AuthorType  AuthorType
	AuthorId    uuid.UUID
	Version     int
	CreatedAt   time.Time
}
