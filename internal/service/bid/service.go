package bid

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"tender-service/internal/mapper"
	"tender-service/internal/model"
	"tender-service/internal/model/dto"
	entity2 "tender-service/internal/model/entity"
	"tender-service/internal/model/entity/bid"
	"tender-service/internal/model/entity/decision"
	"tender-service/internal/model/entity/tender"
	"tender-service/internal/repository"
	service2 "tender-service/internal/service"
	"tender-service/internal/util"
)

type service struct {
	employeeService     service2.EmployeeService
	organizationService service2.OrganizationService
	bidRepository       repository.BidRepository
	tenderService       service2.TenderService
	feedbackRepository  repository.FeedbackRepository
	decisionRepository  repository.DecisionRepository
}

var (
	errStatusCannotBeSelectedByOwner = fmt.Errorf("given status cannot be set by bid owner")
	errCannotVoteOnBid               = fmt.Errorf("cannot vote on given bid")
	errTenderAlreadyClosed           = fmt.Errorf("tender already closed")
)

func NewBidService(
	employeeService service2.EmployeeService,
	organizationService service2.OrganizationService,
	bidRepository repository.BidRepository,
	tenderService service2.TenderService,
	feedbackRepository repository.FeedbackRepository,
	decisionRepository repository.DecisionRepository,
) *service {
	return &service{
		employeeService:     employeeService,
		organizationService: organizationService,
		bidRepository:       bidRepository,
		tenderService:       tenderService,
		feedbackRepository:  feedbackRepository,
		decisionRepository:  decisionRepository,
	}
}

func (s *service) CreateNewBid(ctx context.Context, createDto dto.CreateBidDto) (dto.BidDto, error) {
	if err := s.tenderService.ValidateTenderExists(ctx, createDto.TenderId); err != nil {
		return dto.BidDto{}, err
	}

	if err := s.employeeService.ValidateEmployeeExistsById(ctx, createDto.AuthorId); err != nil {
		return dto.BidDto{}, err
	}

	//newBid := mapper.CreateBidDtoToBid(createDto)

	saved, err := s.bidRepository.SaveBid(ctx, mapper.CreateBidDtoToBid(createDto))
	if err != nil {
		return dto.BidDto{}, err
	}

	return mapper.BidToBidDto(saved), err
}

func (s *service) GetUserBids(ctx context.Context, page util.Page, username string) ([]dto.BidDto, error) {
	user, err := s.employeeService.GetEmployeeByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	bids, err := s.bidRepository.GetBidList(ctx, page, uuid.Nil, user.Id)
	if err != nil {
		return nil, err
	}

	fmt.Println(bids)

	return mapper.BidListToBidDtoList(bids), nil
}

func (s *service) GetTenderBids(ctx context.Context, page util.Page, tenderId uuid.UUID, username string) ([]dto.BidDto, error) {
	if err := s.tenderService.ValidateEmployeeRightsOnTender(ctx, tenderId, username); err != nil {
		return nil, err
	}

	bids, err := s.bidRepository.GetBidList(ctx, page, tenderId, uuid.Nil)
	if err != nil {
		return nil, err
	}

	return mapper.BidListToBidDtoList(bids), nil
}

func (s *service) GetBidStatus(ctx context.Context, bidId uuid.UUID, username string) (bid.Status, error) {
	entity, err := s.bidRepository.GetBidById(ctx, bidId)
	if err != nil {
		return "", err
	}

	if entity.IsVisible() {
		return entity.Status, nil
	}

	if err = s.validateEmployeeRightsOnBid(ctx, bidId, username); err != nil {
		return "", err
	}

	return entity.Status, nil
}

func (s *service) UpdateBidStatus(ctx context.Context, bidId uuid.UUID, username string, status bid.Status) (dto.BidDto, error) {
	err := s.validateEmployeeRightsOnBid(ctx, bidId, username)
	if err != nil {
		return dto.BidDto{}, err
	}

	if !bid.IsSelectableByOwner(status) {
		return dto.BidDto{}, errStatusCannotBeSelectedByOwner
	}

	updated, err := s.bidRepository.UpdateBidStatus(ctx, bidId, status)
	if err != nil {
		return dto.BidDto{}, err
	}

	return mapper.BidToBidDto(updated), err
}

func (s *service) EditBid(ctx context.Context, bidId uuid.UUID, username string, bidDto dto.UpdateBidDto) (dto.BidDto, error) {
	err := s.validateEmployeeRightsOnBid(ctx, bidId, username)
	if err != nil {
		return dto.BidDto{}, err
	}

	updated, err := s.bidRepository.UpdateBid(ctx, bidId, bidDto.Name, bidDto.Description)
	if err != nil {
		return dto.BidDto{}, err
	}

	return mapper.BidToBidDto(updated), err
}

func (s *service) SubmitBidDecision(ctx context.Context, bidId uuid.UUID, username string, verdict decision.Verdict) (dto.BidDto, error) {
	op := "bid_service.submit_bid_decision"
	curBid, err := s.validateEmployeeRightsOnTenderByBid(ctx, bidId, username)
	if err != nil {
		return dto.BidDto{}, err
	}

	if curBid.Status != bid.Published {
		//TODO for each status
		return dto.BidDto{}, model.NewBadRequestError(op, errCannotVoteOnBid)
	}

	ten, err := s.tenderService.GetTenderById(ctx, curBid.TenderId)
	if err != nil {
		return dto.BidDto{}, err
	}

	if ten.Status != tender.Published {
		return dto.BidDto{}, errTenderAlreadyClosed
	}

	fmt.Println("XD" + verdict)

	if verdict == decision.Rejected {
		updatedBid, err := s.bidRepository.UpdateBidStatus(ctx, curBid.Id, bid.Rejected)
		if err != nil {
			return dto.BidDto{}, err
		}
		fmt.Println("gb" + verdict)
		return mapper.BidToBidDto(updatedBid), nil
	}

	fmt.Println("agg" + verdict)

	_, err = s.decisionRepository.SaveDecision(ctx, decision.Decision{
		Verdict:  verdict,
		Username: username,
		BidId:    bidId,
	})
	if err != nil {
		return dto.BidDto{}, err
	}

	approveCount, err := s.decisionRepository.CountDecisionForBid(ctx, bidId)
	if err != nil {
		return dto.BidDto{}, err
	}

	organizationEmployeeCount, err := s.organizationService.GetOrganizationEmployeeCount(ctx, ten.OrganizationId)
	if err != nil {
		return dto.BidDto{}, err
	}

	if approveCount < min(organizationEmployeeCount, 3) {
		return mapper.BidToBidDto(curBid), nil
	}

	updated, err := s.bidRepository.UpdateBidStatus(ctx, bidId, bid.Approved)
	if err != nil {
		return dto.BidDto{}, err
	}

	_, err = s.tenderService.UpdateTenderStatus(ctx, updated.TenderId, username, tender.Closed)
	if err != nil {
		return dto.BidDto{}, err
	}

	return mapper.BidToBidDto(updated), err
}

func (s *service) CreateBidFeedback(ctx context.Context, bidId uuid.UUID, bidFeedback, username string) (dto.BidDto, error) {
	entity, err := s.validateEmployeeRightsOnTenderByBid(ctx, bidId, username)
	if err != nil {
		return dto.BidDto{}, err
	}

	_, err = s.feedbackRepository.SaveFeedback(ctx, entity2.Feedback{
		BidId:       bidId,
		Description: bidFeedback,
		Username:    username,
	})
	if err != nil {
		return dto.BidDto{}, err
	}

	return mapper.BidToBidDto(entity), nil
}

func (s *service) RollbackBid(ctx context.Context, bidId uuid.UUID, username string, version int) (dto.BidDto, error) {
	if err := s.validateEmployeeRightsOnBid(ctx, bidId, username); err != nil {
		return dto.BidDto{}, err
	}

	updated, err := s.bidRepository.RollbackBid(ctx, bidId, version)
	if err != nil {
		return dto.BidDto{}, err
	}

	return mapper.BidToBidDto(updated), err
}

func (s *service) GetBidReviews(ctx context.Context, page util.Page, tenderId uuid.UUID, authorUsername, requesterUsername string) ([]dto.FeedbackDto, error) {
	if err := s.tenderService.ValidateEmployeeRightsOnTender(ctx, tenderId, requesterUsername); err != nil {
		return nil, err
	}

	author, err := s.employeeService.GetEmployeeByUsername(ctx, authorUsername)
	if err != nil {
		return nil, err
	}

	feedback, err := s.feedbackRepository.GetFeedbackListForGroup(ctx, tenderId, author.Id)
	if err != nil {
		return nil, err
	}

	return mapper.FeedbackListToFeedBackDtoList(feedback), nil
}

func (s *service) validateEmployeeRightsOnTenderByBid(ctx context.Context, bidId uuid.UUID, username string) (bid.Bid, error) {
	entity, err := s.bidRepository.GetBidById(ctx, bidId)
	if err != nil {
		return bid.Bid{}, err
	}

	if err = s.tenderService.ValidateEmployeeRightsOnTender(ctx, entity.TenderId, username); err != nil {
		return bid.Bid{}, err
	}
	return entity, nil
}

func (s *service) validateEmployeeRightsOnBid(ctx context.Context, bidId uuid.UUID, username string) error {
	entity, err := s.bidRepository.GetBidById(ctx, bidId)
	if err != nil {
		return err
	}

	if entity.AuthorType == bid.AuthorUser {
		curUser, err := s.employeeService.GetEmployeeByUsername(ctx, username)
		if err != nil {
			return err
		}
		if entity.AuthorId != curUser.Id {
			return fmt.Errorf("not owner")
		}
	} else {
		ok, err := s.organizationService.UsersHasSimilarOrganization(ctx, entity.AuthorId, username)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("no perm")
		}
	}
	return nil
}
