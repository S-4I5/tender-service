package model

type ApiErrorCode string

const (
	NotAuthorizedCode       ApiErrorCode = "not_authorized"
	ForbiddenCode           ApiErrorCode = "forbidden"
	BadRequestCode          ApiErrorCode = "bad_request"
	InternalServerErrorCode ApiErrorCode = "internal_server_error"
	NotFoundCode            ApiErrorCode = "not_found"
	UnprocessableEntityCode ApiErrorCode = "not_processable_entity"
)

type ApiError struct {
	Source string
	Err    error
	Code   ApiErrorCode
}

func (e ApiError) Error() string {
	return e.Source + ":" + string(e.Code) + ":" + e.Err.Error()
}

func NewUnprocessableEntityError(source string, err error) ApiError {
	return newApiError(source, err, UnprocessableEntityCode)
}

func NewNotAuthorizedError(source string, err error) ApiError {
	return newApiError(source, err, NotAuthorizedCode)
}

func NewNotFoundError(source string, err error) ApiError {
	return newApiError(source, err, NotFoundCode)
}

func NewForbiddenError(source string, err error) ApiError {
	return newApiError(source, err, ForbiddenCode)
}

func NewBadRequestError(source string, err error) ApiError {
	return newApiError(source, err, BadRequestCode)
}

func NewInternalServerError(source string, err error) ApiError {
	return newApiError(source, err, InternalServerErrorCode)
}

func newApiError(source string, err error, code ApiErrorCode) ApiError {
	return ApiError{
		Source: source,
		Err:    err,
		Code:   code,
	}
}
