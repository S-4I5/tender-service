package mapper

import (
	"tender-service/internal/model/dto"
	"tender-service/internal/model/entity/bid"
)

func CreateBidDtoToBid(dto dto.CreateBidDto) bid.Bid {
	return bid.Bid{
		Name:        dto.Name,
		Description: dto.Description,
		Status:      dto.Status,
		TenderId:    dto.TenderId,
		AuthorType:  dto.AuthorType,
		AuthorId:    dto.AuthorId,
		Version:     1,
	}
}

func BidToBidDto(entity bid.Bid) dto.BidDto {
	return dto.BidDto{
		Id:          entity.Id,
		Name:        entity.Name,
		Description: entity.Description,
		Status:      entity.Status,
		TenderId:    entity.TenderId,
		AuthorType:  entity.AuthorType,
		AuthorId:    entity.AuthorId,
		Version:     entity.Version,
		CreatedAt:   entity.CreatedAt,
	}
}

func BidListToBidDtoList(list []bid.Bid) []dto.BidDto {
	dtoList := make([]dto.BidDto, len(list))

	for i := 0; i < len(list); i++ {
		dtoList[i] = BidToBidDto(list[i])
	}

	return dtoList
}
