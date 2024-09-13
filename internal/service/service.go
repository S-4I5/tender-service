package service

import (
	"context"
	"github.com/google/uuid"
	"tender-service/internal/model/dto"
	"tender-service/internal/model/entity"
	"tender-service/internal/model/entity/bid"
	"tender-service/internal/model/entity/decision"
	"tender-service/internal/model/entity/tender"
	"tender-service/internal/util"
)

type TenderService interface {
	GetTenders(ctx context.Context, page util.Page, serviceTypes []tender.ServiceType) ([]dto.TenderDto, error)
	CreateNewTender(ctx context.Context, tenderDto dto.CreateTenderDto) (dto.TenderDto, error)
	GetUserTenders(ctx context.Context, page util.Page, username string) ([]dto.TenderDto, error)
	GetTenderStatus(ctx context.Context, tenderId uuid.UUID, username string) (tender.Status, error)
	UpdateTenderStatus(ctx context.Context, tenderId uuid.UUID, username string, status tender.Status) (dto.TenderDto, error)
	EditTender(ctx context.Context, tenderDto dto.UpdateTenderDto, tenderId uuid.UUID, username string) (dto.TenderDto, error)
	RollbackTender(ctx context.Context, tenderId uuid.UUID, username string, version int) (dto.TenderDto, error)
	ValidateTenderExists(ctx context.Context, tenderId uuid.UUID) error
	ValidateEmployeeRightsOnTender(ctx context.Context, tenderId uuid.UUID, username string) error
	GetTenderById(ctx context.Context, tenderId uuid.UUID) (tender.Tender, error)
}

type BidService interface {
	CreateNewBid(ctx context.Context, dto dto.CreateBidDto) (dto.BidDto, error)
	GetUserBids(ctx context.Context, page util.Page, username string) ([]dto.BidDto, error)
	GetTenderBids(ctx context.Context, page util.Page, tenderId uuid.UUID, username string) ([]dto.BidDto, error)
	GetBidStatus(ctx context.Context, bidId uuid.UUID, username string) (bid.Status, error)
	UpdateBidStatus(ctx context.Context, bidId uuid.UUID, username string, status bid.Status) (dto.BidDto, error)
	EditBid(ctx context.Context, bidId uuid.UUID, username string, bidDto dto.UpdateBidDto) (dto.BidDto, error)
	SubmitBidDecision(ctx context.Context, bidId uuid.UUID, username string, verdict decision.Verdict) (dto.BidDto, error)
	CreateBidFeedback(ctx context.Context, bidId uuid.UUID, bidFeedback, username string) (dto.BidDto, error)
	RollbackBid(ctx context.Context, bidId uuid.UUID, username string, version int) (dto.BidDto, error)
	GetBidReviews(ctx context.Context, page util.Page, tenderId uuid.UUID, authorUsername, requesterUsername string) ([]dto.FeedbackDto, error)
}

type OrganizationService interface {
	UsersHasSimilarOrganization(ctx context.Context, userId uuid.UUID, username string) (bool, error)
	ValidateOrganizationExists(ctx context.Context, id uuid.UUID) error
	ValidateEmployeeBelongsToOrganization(ctx context.Context, orgId uuid.UUID, username string) error
	GetOrganizationEmployeeCount(ctx context.Context, id uuid.UUID) (int, error)
	ValidateEmployeeInAnyOrganization(ctx context.Context, userId uuid.UUID) error
}

type EmployeeService interface {
	GetEmployeeByUsername(ctx context.Context, username string) (entity.Employee, error)
	ValidateEmployeeExistsByUsername(ctx context.Context, username string) error
	GetEmployeeByUsernameById(ctx context.Context, id uuid.UUID) (entity.Employee, error)
	ValidateEmployeeExistsById(ctx context.Context, id uuid.UUID) error
}
