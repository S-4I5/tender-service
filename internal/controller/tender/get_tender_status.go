package tender

import (
	"context"
	"encoding/json"
	"net/http"
	"tender-service/internal/model"
)

func (c *controller) GetTenderStatus(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/get_tender_status"
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

		status, err := c.tenderService.GetTenderStatus(ctx, tenderId, username)
		if err != nil {
			c.errHandler.Handler(err, writer)
			return
		}

		if err = json.NewEncoder(writer).Encode(status); err != nil {
			c.errHandler.Handler(model.NewInternalServerError(op, err), writer)
			return
		}
	}
}
