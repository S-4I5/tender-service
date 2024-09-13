package bid

import (
	"context"
	"encoding/json"
	"net/http"
	"tender-service/internal/model"
)

func (c *controller) PutBidFeedback(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/put_bid_feedback"
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

		bidFeedback := request.URL.Query().Get(bidFeedbackQueryParam)
		if bidFeedback == "" {
			c.errHandler.Handler(model.NewBadRequestError(op, errNoBidFeedbackPresented), writer)
			return
		}

		bid, err := c.bidService.CreateBidFeedback(ctx, bidId, bidFeedback, username)
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
