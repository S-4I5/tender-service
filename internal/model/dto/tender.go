package dto

import (
	"github.com/google/uuid"
	"tender-service/internal/model/entity/tender"
)

type CreateTenderDto struct {
	Name            string             `yaml:"name"`
	Description     string             `yaml:"description"`
	Status          tender.Status      `yaml:"status"`
	ServiceType     tender.ServiceType `yaml:"serviceType"`
	OrganizationId  uuid.UUID          `yaml:"organizationId"`
	CreatorUsername string             `yaml:"creatorUsername"`
}

type TenderDto struct {
	Id             uuid.UUID          `yaml:"id"`
	Name           string             `yaml:"name"`
	Description    string             `yaml:"description"`
	Status         tender.Status      `yaml:"status"`
	ServiceType    tender.ServiceType `yaml:"serviceType"`
	OrganizationId uuid.UUID          `yaml:"organizationId"`
	Version        int                `yaml:"version"`
	//CreatorUsername string             `yaml:"creatorUsername"`
}

type UpdateTenderDto struct {
	Name        string             `yaml:"name"`
	Description string             `yaml:"description"`
	ServiceType tender.ServiceType `yaml:"serviceType"`
}
