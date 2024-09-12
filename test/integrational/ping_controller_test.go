package integrational

import (
	"net/http"
	"tender-service/test"
)

func (s *ApiTestSuite) TestPing() {
	actual, err := http.Get(s.host + "/ping")
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer actual.Body.Close()

	test.ValidateResponse(s.T(), actual, "ok", 200)
}
