package mapper

import (
	"tender-service/internal/model/dto"
	"tender-service/internal/model/entity/tender"
)

func CreateTenderDtoToTender(dto dto.CreateTenderDto) tender.Tender {
	return tender.Tender{
		Name:            dto.Name,
		Description:     dto.Description,
		Status:          tender.Created,
		ServiceType:     dto.ServiceType,
		Version:         1,
		OrganizationId:  dto.OrganizationId,
		CreatorUsername: dto.CreatorUsername,
	}
}

func TenderToTenderDto(entity tender.Tender) dto.TenderDto {
	return dto.TenderDto{
		Id:             entity.Id,
		Name:           entity.Name,
		Description:    entity.Description,
		Status:         entity.Status,
		ServiceType:    entity.ServiceType,
		OrganizationId: entity.OrganizationId,
		Version:        entity.Version,
	}
}

func TenderListToTenderDtoList(list []tender.Tender) []dto.TenderDto {
	dtoList := make([]dto.TenderDto, len(list))

	for i := 0; i < len(list); i++ {
		dtoList[i] = TenderToTenderDto(list[i])
	}

	return dtoList
}
