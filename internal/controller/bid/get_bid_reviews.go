package bid

import (
	"context"
	"encoding/json"
	"net/http"
	"tender-service/internal/model"
	"tender-service/internal/util"
)

func (c *controller) GetBidReviews(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "bid_controller/get_bid_reviews"
		writer.Header().Set("Content-Type", "application/json")

		p := util.NewPageFromRequest(request)

		tenderId, err := getTenderIdFromRequest(request)
		if err != nil {
			c.errHandler.Handler(model.NewNotFoundError(op, err), writer)
			return
		}

		authorUsername := request.URL.Query().Get(authorUsernameQueryParam)
		if authorUsername == "" {
			c.errHandler.Handler(model.NewNotAuthorizedError(op, errNoAuthorUsernamePresented), writer)
			return
		}

		requesterUsername := request.URL.Query().Get(requesterUsernameQueryParam)
		if requesterUsername == "" {
			c.errHandler.Handler(model.NewNotAuthorizedError(op, errNoRequesterUsernamePresented), writer)
			return
		}

		reviews, err := c.bidService.GetBidReviews(ctx, p, tenderId, authorUsername, requesterUsername)
		if err != nil {
			c.errHandler.Handler(err, writer)
			return
		}

		if err = json.NewEncoder(writer).Encode(reviews); err != nil {
			c.errHandler.Handler(model.NewInternalServerError(op, err), writer)
			return
		}
	}
}
