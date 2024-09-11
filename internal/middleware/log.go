package middleware

import (
	"fmt"
	"net/http"
)

func GetLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Request:")
		fmt.Println(request)

		next.ServeHTTP(writer, request)

		fmt.Println("Response")
		fmt.Println(writer)
	})
}
