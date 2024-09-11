package tender

import (
	"context"
	"encoding/json"
	"net/http"
	"tender-service/internal/model"
	dto2 "tender-service/internal/model/dto"
)

func (c *controller) PatchTender(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/patch_tender"
		writer.Header().Set("Content-Type", "application/json")

		tenderId, err := getTenderIdFromRequest(request)
		if err != nil {
			c.errHandler.Handler(model.NewNotFoundError(op, err), writer)
			return
		}

		username := request.URL.Query().Get(usernameQueryParam)
		if username == "" {
			c.errHandler.Handler(model.NewNotAuthorizedError(op, errNoUsernameQueryPresented), writer)
			return
		}

		var dto dto2.UpdateTenderDto
		if err := json.NewDecoder(request.Body).Decode(&dto); err != nil {
			c.errHandler.Handler(model.NewUnprocessableEntityError(op, err), writer)
			return
		}

		updated, err := c.tenderService.EditTender(ctx, dto, tenderId, username)
		if err != nil {
			c.errHandler.Handler(err, writer)
			return
		}

		if err = json.NewEncoder(writer).Encode(updated); err != nil {
			c.errHandler.Handler(model.NewInternalServerError(op, err), writer)
			return
		}
	}
}
