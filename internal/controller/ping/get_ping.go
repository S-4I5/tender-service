package ping

import (
	"context"
	"fmt"
	"net/http"
)

func (c *controller) GetPing(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//writer.Header().Set("Content-Type", "text/plain")

		fmt.Fprint(writer, "ok")
	}
}
