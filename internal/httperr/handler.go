package httperr

import (
	"encoding/json"
	"net/http"
	"tender-service/internal/model"
)

type ApiErrorHandler interface {
	Handler(err error, w http.ResponseWriter)
}

type handler struct {
	messageCodes map[model.ApiErrorCode]int
}

func NewApiErrorHandler() *handler {
	return &handler{messageCodes: map[model.ApiErrorCode]int{
		model.NotAuthorizedCode:       401,
		model.InternalServerErrorCode: 500,
		model.ForbiddenCode:           403,
		model.BadRequestCode:          400,
		model.NotFoundCode:            404,
		model.UnprocessableEntityCode: 422,
	}}
}

func (h *handler) Handler(err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError

	apiErr, ok := err.(model.ApiError)
	if ok {
		status = h.messageCodes[apiErr.Code]
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorDto{Reason: err.Error()})
}
