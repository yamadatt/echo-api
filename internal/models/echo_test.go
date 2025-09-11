package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewEchoRequest(t *testing.T) {
	method := "GET"
	path := "/test"
	headers := map[string]string{"Content-Type": "application/json"}
	queryParams := map[string]string{"param1": "value1"}
	body := "test body"

	req := NewEchoRequest(method, path, headers, queryParams, body)

	if req.Method != method {
		t.Errorf("Expected method %s, got %s", method, req.Method)
	}
	if req.Path != path {
		t.Errorf("Expected path %s, got %s", path, req.Path)
	}
	if req.Headers["Content-Type"] != "application/json" {
		t.Errorf("Expected header Content-Type: application/json, got %s", req.Headers["Content-Type"])
	}
	if req.QueryParams["param1"] != "value1" {
		t.Errorf("Expected query param param1: value1, got %s", req.QueryParams["param1"])
	}
	if req.Body != body {
		t.Errorf("Expected body %s, got %s", body, req.Body)
	}
	if req.Timestamp == "" {
		t.Error("Expected timestamp to be set")
	}
}

func TestNewEchoResponse(t *testing.T) {
	req := &EchoRequest{
		Method:      "GET",
		Path:        "/test",
		Headers:     map[string]string{},
		QueryParams: map[string]string{},
		Body:        "",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	message := "Test message"

	resp := NewEchoResponse(req, message)

	if resp.Request.Method != req.Method {
		t.Errorf("Expected request method %s, got %s", req.Method, resp.Request.Method)
	}
	if resp.Message != message {
		t.Errorf("Expected message %s, got %s", message, resp.Message)
	}
	if resp.ProcessedAt == "" {
		t.Error("Expected processedAt to be set")
	}
}

func TestNewErrorResponse(t *testing.T) {
	errorType := "Test Error"
	message := "Test error message"

	resp := NewErrorResponse(errorType, message)

	if resp.Error != errorType {
		t.Errorf("Expected error %s, got %s", errorType, resp.Error)
	}
	if resp.Message != message {
		t.Errorf("Expected message %s, got %s", message, resp.Message)
	}
	if resp.Timestamp == "" {
		t.Error("Expected timestamp to be set")
	}
}

func TestEchoResponseToJSON(t *testing.T) {
	req := &EchoRequest{
		Method:      "GET",
		Path:        "/test",
		Headers:     map[string]string{"Accept": "application/json"},
		QueryParams: map[string]string{"test": "value"},
		Body:        "",
		Timestamp:   "2023-01-01T00:00:00Z",
	}
	
	resp := &EchoResponse{
		Request:     *req,
		Message:     "Test response",
		ProcessedAt: "2023-01-01T00:00:01Z",
	}

	jsonStr, err := resp.ToJSON()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Parse the JSON back to verify structure
	var parsed EchoResponse
	err = json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed.Request.Method != req.Method {
		t.Errorf("Expected method %s, got %s", req.Method, parsed.Request.Method)
	}
	if parsed.Message != resp.Message {
		t.Errorf("Expected message %s, got %s", resp.Message, parsed.Message)
	}
}

func TestErrorResponseToJSON(t *testing.T) {
	resp := &ErrorResponse{
		Error:     "Test Error",
		Message:   "Test error message",
		Timestamp: "2023-01-01T00:00:00Z",
	}

	jsonStr, err := resp.ToJSON()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Parse the JSON back to verify structure
	var parsed ErrorResponse
	err = json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed.Error != resp.Error {
		t.Errorf("Expected error %s, got %s", resp.Error, parsed.Error)
	}
	if parsed.Message != resp.Message {
		t.Errorf("Expected message %s, got %s", resp.Message, parsed.Message)
	}
}