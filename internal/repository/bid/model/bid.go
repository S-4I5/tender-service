package model

import (
	"github.com/google/uuid"
	"tender-service/internal/model/entity/bid"
	"time"
)

type Bid struct {
	Id           uuid.UUID
	Status       string
	Decision     bid.Decision
	TenderId     uuid.UUID
	AuthorType   string
	BidVersionId uuid.UUID
	AuthorId     uuid.UUID
	CreatedAt    time.Time
}

type BidVersion struct {
	Id          uuid.UUID
	BidId       uuid.UUID
	Name        string
	Description string
	Version     int
}

type BidSum struct {
	Id          uuid.UUID
	Name        string
	Description string
	Decision    bid.Decision
	Version     int
	Status      string
	TenderId    uuid.UUID
	AuthorType  string
	AuthorId    uuid.UUID
	CreatedAt   time.Time
}

func MergeBidAndVersionToBid(v BidVersion, b Bid) bid.Bid {
	return bid.Bid{
		Id:          b.Id,
		Name:        v.Name,
		Description: v.Description,
		Decision:    b.Decision,
		Status:      bid.Status(b.Status),
		TenderId:    b.TenderId,
		AuthorType:  bid.AuthorType(b.AuthorType),
		AuthorId:    b.AuthorId,
		Version:     v.Version,
		CreatedAt:   b.CreatedAt,
	}
}

func BidSumToBid(sum BidSum) bid.Bid {
	return bid.Bid{
		Id:          sum.Id,
		Name:        sum.Name,
		Decision:    sum.Decision,
		Description: sum.Description,
		Status:      bid.Status(sum.Status),
		TenderId:    sum.TenderId,
		AuthorType:  bid.AuthorType(sum.AuthorType),
		AuthorId:    sum.AuthorId,
		Version:     sum.Version,
		CreatedAt:   sum.CreatedAt,
	}
}

func BidSumListToBidList(list []BidSum) []bid.Bid {
	dtoList := make([]bid.Bid, len(list))

	for i := 0; i < len(list); i++ {
		dtoList[i] = BidSumToBid(list[i])
	}

	return dtoList
}
