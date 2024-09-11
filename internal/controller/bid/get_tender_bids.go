package bid

import (
	"context"
	"encoding/json"
	"net/http"
	"tender-service/internal/model"
	"tender-service/internal/util"
)

func (c *controller) GetTenderBids(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/get_tender_bids"
		writer.Header().Set("Content-Type", "application/json")

		p := util.NewPageFromRequest(request)

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

		bids, err := c.bidService.GetTenderBids(ctx, p, tenderId, username)
		if err != nil {
			c.errHandler.Handler(err, writer)
			return
		}

		if err = json.NewEncoder(writer).Encode(bids); err != nil {
			c.errHandler.Handler(model.NewInternalServerError(op, err), writer)
			return
		}
	}
}
