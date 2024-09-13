package ping

import (
	"context"
	"fmt"
	"net/http"
)

func (c *controller) GetPing(ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "ok")
	}
}
