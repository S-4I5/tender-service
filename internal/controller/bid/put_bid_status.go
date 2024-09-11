package bid

import (
	"context"
	"encoding/json"
	"net/http"
	"tender-service/internal/model"
	"tender-service/internal/model/entity/bid"
)

func (c *controller) PutBidStatus(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/put_bid_status"
		writer.Header().Set("Content-Type", "application/json")

		bidId, err := getBidIdFromRequest(request)
		if err != nil {
			c.errHandler.Handler(model.NewNotFoundError(op, err), writer)
			return
		}

		username := request.URL.Query().Get(usernameQueryParam)
		if username == "" {
			c.errHandler.Handler(model.NewNotAuthorizedError(op, errNoUsernameQueryPresented), writer)
			return
		}

		status := bid.Status(request.URL.Query().Get(statusQueryParam))

		updated, err := c.bidService.UpdateBidStatus(ctx, bidId, username, status)
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
