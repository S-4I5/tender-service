package repository

import (
	"context"
	"github.com/google/uuid"
	"tender-service/internal/model/entity"
	"tender-service/internal/model/entity/bid"
	"tender-service/internal/model/entity/decision"
	"tender-service/internal/model/entity/organization"
	"tender-service/internal/model/entity/tender"
	"tender-service/internal/util"
)

type EmployeeRepository interface {
	GetEmployeeByUsername(ctx context.Context, username string) (entity.Employee, error)
	EmployeeExistByUsername(ctx context.Context, username string) (bool, error)
	GetEmployeeById(ctx context.Context, id uuid.UUID) (entity.Employee, error)
	EmployeeExistById(ctx context.Context, id uuid.UUID) (bool, error)
}

type OrganizationRepository interface {
	GetOrganizationById(ctx context.Context, id uuid.UUID) (organization.Organization, error)
	OrganizationExistById(ctx context.Context, id uuid.UUID) (bool, error)
}

type OrganizationResponsibleRepository interface {
	IsResponsibleInOrganization(ctx context.Context, username string, organizationId uuid.UUID) (bool, error)
	CountEmployeesInOrganization(ctx context.Context, organizationId uuid.UUID) (int, error)
	UsersHasSimilarOrganization(ctx context.Context, userId uuid.UUID, username string) (bool, error)
}

type TenderRepository interface {
	SaveTender(ctx context.Context, version tender.Tender) (tender.Tender, error)
	GetTenderById(ctx context.Context, id uuid.UUID) (tender.Tender, error)
	GetTenderList(ctx context.Context, page util.Page, serviceTypes []tender.ServiceType, username string, onlyPublished bool) ([]tender.Tender, error)
	UpdateTender(ctx context.Context, id uuid.UUID, name, description string, serviceType tender.ServiceType) (tender.Tender, error)
	UpdateTenderStatus(ctx context.Context, id uuid.UUID, status tender.Status) (tender.Tender, error)
	RollbackTender(ctx context.Context, id uuid.UUID, version int) (tender.Tender, error)
}

type BidRepository interface {
	SaveBid(ctx context.Context, version bid.Bid) (bid.Bid, error)
	GetBidById(ctx context.Context, id uuid.UUID) (bid.Bid, error)
	GetBidList(ctx context.Context, page util.Page, tenderId uuid.UUID, userId uuid.UUID) ([]bid.Bid, error)
	UpdateBidStatus(ctx context.Context, id uuid.UUID, stat bid.Status) (bid.Bid, error)
	UpdateBid(ctx context.Context, id uuid.UUID, name, description string) (bid.Bid, error)
	RollbackBid(ctx context.Context, id uuid.UUID, version int) (bid.Bid, error)
}

type DecisionRepository interface {
	SaveDecision(ctx context.Context, decision decision.Decision) (decision.Decision, error)
	CountDecisionForBid(ctx context.Context, bidId uuid.UUID) (int, error)
}

type FeedbackRepository interface {
	SaveFeedback(ctx context.Context, feedback entity.Feedback) (entity.Feedback, error)
	GetFeedbackListForGroup(ctx context.Context, tenderId uuid.UUID, userId uuid.UUID) ([]entity.Feedback, error)
}
