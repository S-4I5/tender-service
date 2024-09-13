package dto

import (
	"github.com/google/uuid"
	"tender-service/internal/model/entity/tender"
)

type CreateTenderDto struct {
	Name            string             `json:"name" validate:"required"`
	Description     string             `json:"description" validate:"required"`
	ServiceType     tender.ServiceType `json:"serviceType" validate:"required"`
	OrganizationId  uuid.UUID          `json:"organizationId" validate:"required"`
	CreatorUsername string             `json:"creatorUsername" validate:"required"`
}

type TenderDto struct {
	Id             uuid.UUID          `json:"id"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	Status         tender.Status      `json:"status"`
	ServiceType    tender.ServiceType `json:"serviceType"`
	OrganizationId uuid.UUID          `json:"organizationId"`
	Version        int                `json:"version"`
}

type UpdateTenderDto struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	ServiceType tender.ServiceType `json:"serviceType"`
}
