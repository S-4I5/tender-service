package model

import (
	"database/sql"
	"github.com/google/uuid"
	"tender-service/internal/model/entity/organization"
	"time"
)

type Organization struct {
	Id          uuid.UUID      `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Type        sql.NullString `db:"type"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func DbOrganizationToOrganization(org Organization) organization.Organization {
	return organization.Organization{
		Id:          org.Id,
		Name:        org.Name,
		Description: org.Description.String,
		Type:        organization.Type(org.Type.String),
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}
}
