package test

import (
	"bytes"
	"encoding/json"
	"github.com/nsf/jsondiff"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func HttpPut(url string, dto any) (*http.Response, error) {
	return do(http.MethodPut, url, dto)
}

func HttpPatch(url string, dto any) (*http.Response, error) {
	return do(http.MethodPatch, url, dto)
}

func do(method, url string, dto any) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, ToBuffer(dto))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ValidateJsonStringResponse(t *testing.T, resp *http.Response, expected string, code int) {
	ValidateResponse(t, resp, "\""+expected+"\"\n", code)
}

func ValidateResponse(t *testing.T, resp *http.Response, expected string, code int) {
	if resp.StatusCode != code {
		t.Errorf("Expected status %d, got %v", code, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != expected {
		t.Errorf("Unexpected response body: got %v, want %v", string(body), expected)
	}
}

func ValidateJsonResponse(t *testing.T, resp *http.Response, expected []byte, code int) {
	if resp.StatusCode != code {
		t.Errorf("Expected status %d, got %v", code, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if !CompareJson(expected, body) {
		t.Errorf("Response body don't match: %v", err)
	}
}

func ToBuffer(o any) *bytes.Buffer {
	result, _ := json.Marshal(o)
	return bytes.NewBuffer(result)
}

func ReadJson(path string) []byte {
	filePath := "../resources"

	file, err := os.Open(filePath + path + ".json")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	return content
}

func CompareJson(expected, given []byte) bool {
	result, _ := jsondiff.Compare(given, expected, &jsondiff.Options{})
	return result == jsondiff.SupersetMatch || result == jsondiff.FullMatch
}
