package tender

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"tender-service/internal/httperr"
	"tender-service/internal/service"
)

type controller struct {
	errHandler    httperr.ApiErrorHandler
	tenderService service.TenderService
	validator     *validator.Validate
}

const (
	usernameQueryParam    = "username"
	tenderIdPathValue     = "tenderId"
	serviceTypeQueryParam = "service_type"
	versionPathValue      = "version"
	statusQueryParam      = "status"
)

var (
	errTenderPathValueNotFound  = fmt.Errorf("path value tenderId is not presented")
	errNoUsernameQueryPresented = fmt.Errorf("request param username is not presented")
	errIncorrectServiceType     = fmt.Errorf("provided incorrect service type")
	errIncorrectTenderStatus    = fmt.Errorf("incorrect tender status")
)

func NewTenderController(tenderService service.TenderService, errHandler httperr.ApiErrorHandler) *controller {
	return &controller{
		tenderService: tenderService,
		errHandler:    errHandler,
		validator:     validator.New(validator.WithRequiredStructEnabled()),
	}
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
