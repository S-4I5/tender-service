package tender

import (
	"context"
	"encoding/json"
	"net/http"
	"tender-service/internal/model"
	"tender-service/internal/util"
)

func (c *controller) GetUserTenders(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/get_user_tenders"
		writer.Header().Set("Content-Type", "application/json")

		p := util.NewPageFromRequest(request)

		username := request.URL.Query().Get(usernameQueryParam)
		if username == "" {
			c.errHandler.Handler(model.NewNotAuthorizedError(op, errNoUsernameQueryPresented), writer)
			return
		}

		tenders, err := c.tenderService.GetUserTenders(ctx, p, username)
		if err != nil {
			c.errHandler.Handler(err, writer)
			return
		}

		if err = json.NewEncoder(writer).Encode(tenders); err != nil {
			c.errHandler.Handler(model.NewInternalServerError(op, err), writer)
			return
		}
	}
}
