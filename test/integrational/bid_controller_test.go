package integrational

import (
	"io"
	"net/http"
)

func (s *ApiTestSuite) sTestCreateBid() {
	resp, err := http.Get(s.host + "/te")
	if err != nil {
		s.T().Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.T().Errorf("Expected status 200 OK, got %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.T().Fatalf("Failed to read response body: %v", err)
	}

	expected := "\"ok\"\n"
	if string(body) != expected {
		s.T().Errorf("Unexpected response body: got %v, want %v", string(body), expected)
	}
}
