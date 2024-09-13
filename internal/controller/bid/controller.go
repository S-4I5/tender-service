package bid

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"tender-service/internal/httperr"
	"tender-service/internal/service"
)

type controller struct {
	bidService service.BidService
	errHandler httperr.ApiErrorHandler
	validator  *validator.Validate
}

const (
	usernameQueryParam          = "username"
	tenderIdPathValue           = "tenderId"
	bidIdPathValue              = "bidId"
	versionPathValue            = "version"
	statusQueryParam            = "status"
	decisionQueryParam          = "decision"
	bidFeedbackQueryParam       = "bidFeedback"
	authorUsernameQueryParam    = "authorUsername"
	requesterUsernameQueryParam = "requesterUsername"
)

var (
	errTenderPathValueNotFound      = fmt.Errorf("path value tenderId is not presented")
	errBidPathValueNotFound         = fmt.Errorf("path value bidId is not presented")
	errNoUsernameQueryPresented     = fmt.Errorf("request param username is not presented")
	errNoAuthorUsernamePresented    = fmt.Errorf("request param authorUsername is not presented")
	errNoRequesterUsernamePresented = fmt.Errorf("request param requesterUsername is not presented")
	errNoBidFeedbackPresented       = fmt.Errorf("request param bidFeedback is not presented")
	errIncorrectBidDecision         = fmt.Errorf("incorrect bid decision")
)

func NewBidController(bidService service.BidService, errHandler httperr.ApiErrorHandler) *controller {
	return &controller{
		bidService: bidService,
		errHandler: errHandler,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
	}
}

func getBidIdFromRequest(request *http.Request) (uuid.UUID, error) {
	tenderId := request.PathValue(bidIdPathValue)
	if tenderId == "" {
		return uuid.Nil, errBidPathValueNotFound
	}
	tenderUuid, err := uuid.Parse(tenderId)
	if err != nil {
		return uuid.Nil, err
	}
	return tenderUuid, nil
}

func getTenderIdFromRequest(request *http.Request) (uuid.UUID, error) {
	tenderId := request.PathValue(tenderIdPathValue)
	if tenderId == "" {
		return uuid.Nil, errTenderPathValueNotFound
	}
	tenderUuid, err := uuid.Parse(tenderId)
	if err != nil {
		return uuid.Nil, err
	}
	return tenderUuid, nil
}
