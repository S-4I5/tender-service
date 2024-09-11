package tender

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	Closed    Status = "Closed"
	Created   Status = "Created"
	Published Status = "Published"
)

func IsTenderStatus(status string) bool {
	mapped := Status(status)
	return mapped == Closed || mapped == Created || mapped == Published
}

type ServiceType string

const (
	Construction ServiceType = "Construction"
	Delivery     ServiceType = "Delivery"
	Manufacture  ServiceType = "Manufacture"
)

func IsServiceType(serviceType string) bool {
	mapped := ServiceType(serviceType)
	return mapped == Construction || mapped == Delivery || mapped == Manufacture
}

type Tender struct {
	Id              uuid.UUID
	Name            string
	Description     string
	Status          Status
	ServiceType     ServiceType
	Version         int
	CreatedAt       time.Time
	OrganizationId  uuid.UUID
	CreatorUsername string
}
