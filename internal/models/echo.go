package models

import (
	"encoding/json"
	"time"
)

// EchoRequest represents the incoming HTTP request data
type EchoRequest struct {
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"queryParams"`
	Body        string            `json:"body,omitempty"`
	Timestamp   string            `json:"timestamp"`
}

// EchoResponse represents the response containing the echo of the request
type EchoResponse struct {
	Request     EchoRequest `json:"request"`
	Message     string      `json:"message"`
	ProcessedAt string      `json:"processedAt"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// NewEchoRequest creates a new EchoRequest with current timestamp
func NewEchoRequest(method, path string, headers, queryParams map[string]string, body string) *EchoRequest {
	return &EchoRequest{
		Method:      method,
		Path:        path,
		Headers:     headers,
		QueryParams: queryParams,
		Body:        body,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// NewEchoResponse creates a new EchoResponse with current processed timestamp
func NewEchoResponse(request *EchoRequest, message string) *EchoResponse {
	return &EchoResponse{
		Request:     *request,
		Message:     message,
		ProcessedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewErrorResponse creates a new ErrorResponse with current timestamp
func NewErrorResponse(error, message string) *ErrorResponse {
	return &ErrorResponse{
		Error:     error,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// ToJSON converts the EchoResponse to JSON string
func (e *EchoResponse) ToJSON() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts the ErrorResponse to JSON string
func (e *ErrorResponse) ToJSON() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data), nil
}