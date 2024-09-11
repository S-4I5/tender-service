package mapper

import (
	"tender-service/internal/model/dto"
	"tender-service/internal/model/entity"
)

func FeedbackToFeedBackDto(feedback entity.Feedback) dto.FeedbackDto {
	return dto.FeedbackDto{
		Id:          feedback.Id,
		Description: feedback.Description,
		CreatedAt:   feedback.CreatedAt,
	}
}

func FeedbackListToFeedBackDtoList(list []entity.Feedback) []dto.FeedbackDto {
	dtoList := make([]dto.FeedbackDto, len(list))

	for i := 0; i < len(list); i++ {
		dtoList[i] = FeedbackToFeedBackDto(list[i])
	}

	return dtoList
}
