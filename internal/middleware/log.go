package middleware

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s at %s", r.Method, r.URL.Path, start.String())

		if r.Body != nil {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			log.Printf("Request Body: %s", string(bodyBytes))

			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		rw := &responseWriter{w, http.StatusOK, &bytes.Buffer{}}

		next.ServeHTTP(rw, r)

		log.Printf("Response Body: %s", rw.body.String())
		log.Printf("Completed %s %s with %d in %v", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
