package tender

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"tender-service/internal/model"
)

func (c *controller) PutTenderRollback(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		op := "tender_controller/put_tender_rollback"
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

		versionString := request.PathValue(versionPathValue)
		version, err := strconv.Atoi(versionString)
		if err != nil {
			c.errHandler.Handler(model.NewBadRequestError(op, err), writer)
			return
		}

		tender, err := c.tenderService.RollbackTender(ctx, tenderId, username, version)
		if err != nil {
			c.errHandler.Handler(err, writer)
			return
		}

		if err = json.NewEncoder(writer).Encode(tender); err != nil {
			c.errHandler.Handler(model.NewInternalServerError(op, err), writer)
			return
		}
	}
}
