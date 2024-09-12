package tender

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"tender-service/internal/model"
	"tender-service/internal/model/entity/tender"
	"tender-service/internal/util"
)

func (c *controller) GetTenders(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/get_tenders"
		writer.Header().Set("Content-Type", "application/json")

		p := util.NewPageFromRequest(request)

		rowServiceTypesString := request.URL.Query().Get(serviceTypeQueryParam)
		var serviceTypes []tender.ServiceType

		if rowServiceTypesString != "" {
			log.Println("in", rowServiceTypesString)
			rowServiceTypes := strings.Split(rowServiceTypesString, ",")
			log.Println("rows", rowServiceTypes)
			serviceTypes = make([]tender.ServiceType, len(rowServiceTypes))
			for i := 0; i < len(rowServiceTypes); i++ {
				log.Println("cur", serviceTypes[i])
				if !tender.IsServiceType(rowServiceTypes[i]) {
					c.errHandler.Handler(model.NewBadRequestError(op, errIncorrectServiceType), writer)
					return
				}
				serviceTypes[i] = tender.ServiceType(rowServiceTypes[i])
			}
		}

		log.Println("type", serviceTypes)

		tenders, err := c.tenderService.GetTenders(ctx, p, serviceTypes)
		if err != nil {
			c.errHandler.Handler(err, writer)
			return
		}

		if err = json.NewEncoder(writer).Encode(tenders); err != nil {
			c.errHandler.Handler(model.NewInternalServerError(op, err), writer)
			return
		}
	}
}
