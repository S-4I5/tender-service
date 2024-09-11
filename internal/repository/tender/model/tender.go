package model

import (
	"database/sql"
	"github.com/google/uuid"
	"tender-service/internal/model/entity/tender"
	"time"
)

type Tender struct {
	Id              uuid.UUID
	Status          string
	TenderVersionId sql.NullInt32
	OrganizationId  uuid.UUID
	CreatorUsername string
	CreatedAt       time.Time
}

type TenderVersion struct {
	Id          int
	TenderId    uuid.UUID
	Name        string
	Description string
	ServiceType string
	Version     int
}

type TenderSum struct {
	Id              uuid.UUID
	Status          string
	Name            string
	Description     string
	ServiceType     string
	Version         int
	OrganizationId  uuid.UUID
	CreatorUsername string
	CreatedAt       time.Time
}

func DbTenderSumToTender(tenderSum TenderSum) tender.Tender {
	return tender.Tender{
		Id:              tenderSum.Id,
		Name:            tenderSum.Name,
		Description:     tenderSum.Description,
		Status:          tender.Status(tenderSum.Status),
		ServiceType:     tender.ServiceType(tenderSum.ServiceType),
		Version:         tenderSum.Version,
		CreatedAt:       tenderSum.CreatedAt,
		OrganizationId:  tenderSum.OrganizationId,
		CreatorUsername: tenderSum.CreatorUsername,
	}
}

func MergeTenderWithVersion(v TenderVersion, t Tender) TenderSum {
	return TenderSum{
		Id:              t.Id,
		Status:          t.Status,
		Name:            v.Name,
		Description:     v.Description,
		ServiceType:     v.ServiceType,
		Version:         v.Version,
		OrganizationId:  t.OrganizationId,
		CreatorUsername: t.CreatorUsername,
		CreatedAt:       t.CreatedAt,
	}
}

func DdTenderVersionListToTenderList(list []TenderSum) []tender.Tender {
	dtoList := make([]tender.Tender, len(list))

	for i := 0; i < len(list); i++ {
		dtoList[i] = DbTenderSumToTender(list[i])
	}

	return dtoList
}
