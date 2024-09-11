package ping

import (
	"context"
	"encoding/json"
	"net/http"
)

func (c *controller) GetPing(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")

		if err := json.NewEncoder(writer).Encode("ok"); err != nil {
			return
		}
	}
}
