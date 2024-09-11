package util

import (
	"net/http"
	"strconv"
)

type Page struct {
	Offset int
	Limit  int
}

const (
	offsetQueryParam = "offset"
	limitQueryParam  = "limit"
)

func NewPage(offset, limit int) Page {
	return Page{
		Offset: 0,
		Limit:  0,
	}
}

func NewPageFromRequest(request *http.Request) Page {
	return Page{
		Offset: getUrlRequestParamOrDefault(offsetQueryParam, 0, request),
		Limit:  getUrlRequestParamOrDefault(limitQueryParam, 5, request),
	}
}

func getUrlRequestParamOrDefault(requestParamName string, def int, request *http.Request) int {
	requestParam := request.URL.Query().Get(requestParamName)
	value, err := strconv.Atoi(requestParam)
	if err != nil {
		return def
	}
	return value
}
