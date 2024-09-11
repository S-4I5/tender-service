package tender

import (
	"context"
	"encoding/json"
	"net/http"
	"tender-service/internal/model"
	"tender-service/internal/model/entity/tender"
)

func (c *controller) PutTenderStatus(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/put_tender_status"
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

		status := request.URL.Query().Get(statusQueryParam)
		if ok := tender.IsTenderStatus(status); !ok {
			c.errHandler.Handler(model.NewBadRequestError(op, errIncorrectTenderStatus), writer)
			return
		}

		updated, err := c.tenderService.UpdateTenderStatus(ctx, tenderId, username, tender.Status(status))
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
