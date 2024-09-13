package bid

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"tender-service/internal/model"
	"tender-service/internal/model/entity/decision"
)

func (c *controller) PutBidSubmitDecision(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/put_bid_submit_decision"
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

		des := request.URL.Query().Get(decisionQueryParam)
		if !decision.IsDecisionVerdict(des) {
			c.errHandler.Handler(model.NewNotAuthorizedError(op, errIncorrectBidDecision), writer)
			return
		}

		log.Println(des)

		bid, err := c.bidService.SubmitBidDecision(ctx, bidId, username, decision.Verdict(des))
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
