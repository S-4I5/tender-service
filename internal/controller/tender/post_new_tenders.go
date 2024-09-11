package tender

import (
	"context"
	"encoding/json"
	"net/http"
	"tender-service/internal/model"
	dto2 "tender-service/internal/model/dto"
)

func (c *controller) PostNewTender(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/post_new_tender"
		writer.Header().Set("Content-Type", "application/json")

		var dto dto2.CreateTenderDto
		if err := json.NewDecoder(request.Body).Decode(&dto); err != nil {
			c.errHandler.Handler(model.NewUnprocessableEntityError(op, err), writer)
			return
		}

		if err := c.validator.Struct(dto); err != nil {
			c.errHandler.Handler(model.NewBadRequestError(op, err), writer)
			return
		}

		saved, err := c.tenderService.CreateNewTender(ctx, dto)
		if err != nil {
			c.errHandler.Handler(err, writer)
			return
		}

		if err = json.NewEncoder(writer).Encode(saved); err != nil {
			c.errHandler.Handler(model.NewInternalServerError(op, err), writer)
			return
		}
	}
}
