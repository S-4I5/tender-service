package bid

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"tender-service/internal/model"
)

func (c *controller) PutBidRollback(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/put_bid_rollback"
		writer.Header().Set("Content-Type", "application/json")

		bidId, err := getBidIdFromRequest(request)
		if err != nil {
			c.errHandler.Handler(model.NewNotFoundError(op, err), writer)
			return
		}

		versionString := request.PathValue(versionPathValue)
		version, err := strconv.Atoi(versionString)
		if err != nil {
			c.errHandler.Handler(model.NewBadRequestError(op, err), writer)
			return
		}

		username := request.URL.Query().Get(usernameQueryParam)
		if username == "" {
			c.errHandler.Handler(model.NewNotAuthorizedError(op, errNoUsernameQueryPresented), writer)
			return
		}

		bid, err := c.bidService.RollbackBid(ctx, bidId, username, version)
		if err != nil {
			c.errHandler.Handler(err, writer)
			return
		}

		if err = json.NewEncoder(writer).Encode(bid); err != nil {
			c.errHandler.Handler(model.NewInternalServerError(op, err), writer)
			return
		}
	}
}
