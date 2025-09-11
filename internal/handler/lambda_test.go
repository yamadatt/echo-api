package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"echo-api/internal/models"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleRequest_GET(t *testing.T) {
	handler := NewLambdaHandler()
	ctx := context.Background()

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/test",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   "test-agent",
		},
		QueryStringParameters: map[string]string{
			"param1": "value1",
			"param2": "value2",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			Stage: "test",
		},
	}

	response, err := handler.HandleRequest(ctx, request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Parse response body
	var echoResponse models.EchoResponse
	err = json.Unmarshal([]byte(response.Body), &echoResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Verify request data
	if echoResponse.Request.Method != "GET" {
		t.Errorf("Expected method GET, got %s", echoResponse.Request.Method)
	}
	if echoResponse.Request.Path != "/test" {
		t.Errorf("Expected path /test, got %s", echoResponse.Request.Path)
	}
	if echoResponse.Request.Headers["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type header, got %s", echoResponse.Request.Headers["Content-Type"])
	}
	if echoResponse.Request.QueryParams["param1"] != "value1" {
		t.Errorf("Expected query param param1=value1, got %s", echoResponse.Request.QueryParams["param1"])
	}
}

func TestHandleRequest_POST(t *testing.T) {
	handler := NewLambdaHandler()
	ctx := context.Background()

	requestBody := `{"key": "value", "number": 123}`
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/api/echo",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: requestBody,
		RequestContext: events.APIGatewayProxyRequestContext{
			Stage: "test",
		},
	}

	response, err := handler.HandleRequest(ctx, request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Parse response body
	var echoResponse models.EchoResponse
	err = json.Unmarshal([]byte(response.Body), &echoResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Verify request data
	if echoResponse.Request.Method != "POST" {
		t.Errorf("Expected method POST, got %s", echoResponse.Request.Method)
	}
	if echoResponse.Request.Body != requestBody {
		t.Errorf("Expected body %s, got %s", requestBody, echoResponse.Request.Body)
	}
}

func TestHandleRequest_MethodNotAllowed(t *testing.T) {
	handler := NewLambdaHandler()
	ctx := context.Background()

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "DELETE",
		Path:       "/test",
		RequestContext: events.APIGatewayProxyRequestContext{
			Stage: "test",
		},
	}

	response, err := handler.HandleRequest(ctx, request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, response.StatusCode)
	}

	// Parse error response
	var errorResponse models.ErrorResponse
	err = json.Unmarshal([]byte(response.Body), &errorResponse)
	if err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}

	if errorResponse.Error != "Method Not Allowed" {
		t.Errorf("Expected error 'Method Not Allowed', got %s", errorResponse.Error)
	}
}

func TestIsMethodAllowed(t *testing.T) {
	handler := NewLambdaHandler()

	testCases := []struct {
		method   string
		expected bool
	}{
		{"GET", true},
		{"POST", true},
		{"OPTIONS", true},
		{"PUT", false},
		{"DELETE", false},
		{"PATCH", false},
		{"HEAD", false},
	}

	for _, tc := range testCases {
		result := handler.isMethodAllowed(tc.method)
		if result != tc.expected {
			t.Errorf("For method %s, expected %v, got %v", tc.method, tc.expected, result)
		}
	}
}

func TestParseRequest(t *testing.T) {
	handler := NewLambdaHandler()

	request := &events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/api/test",
		Headers: map[string]string{
			"Content-Type":   "application/json",
			"Authorization": "Bearer token123",
		},
		QueryStringParameters: map[string]string{
			"filter": "active",
			"limit":  "10",
		},
		Body: `{"data": "test"}`,
	}

	echoRequest := handler.parseRequest(request)

	if echoRequest.Method != "POST" {
		t.Errorf("Expected method POST, got %s", echoRequest.Method)
	}
	if echoRequest.Path != "/api/test" {
		t.Errorf("Expected path /api/test, got %s", echoRequest.Path)
	}
	if echoRequest.Headers["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type header, got %s", echoRequest.Headers["Content-Type"])
	}
	if echoRequest.QueryParams["filter"] != "active" {
		t.Errorf("Expected query param filter=active, got %s", echoRequest.QueryParams["filter"])
	}
	if echoRequest.Body != `{"data": "test"}` {
		t.Errorf("Expected body to match, got %s", echoRequest.Body)
	}
	if echoRequest.Timestamp == "" {
		t.Error("Expected timestamp to be set")
	}
}