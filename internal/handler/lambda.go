package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"echo-api/internal/models"
	"echo-api/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
)

// LambdaHandler handles AWS Lambda proxy requests
type LambdaHandler struct {
	logger *logger.Logger
}

// NewLambdaHandler creates a new Lambda handler instance
func NewLambdaHandler() *LambdaHandler {
	return &LambdaHandler{
		logger: logger.New(),
	}
}

// HandleRequest processes the incoming API Gateway proxy request
func (h *LambdaHandler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	h.logger.Info("Processing request", map[string]interface{}{
		"method":       request.HTTPMethod,
		"path":         request.Path,
		"stage":        request.RequestContext.Stage,
		"resource":     request.Resource,
		"full_request": fmt.Sprintf("%+v", request),
	})

	// Check if method is allowed (GET or POST)
	if !h.isMethodAllowed(request.HTTPMethod) {
		h.logger.Warn("Method not allowed", map[string]interface{}{
			"method": request.HTTPMethod,
		})
		return h.createErrorResponse(http.StatusMethodNotAllowed, "Method Not Allowed", "Only GET and POST methods are supported")
	}

	// Parse the request
	echoRequest := h.parseRequest(&request)

	// Create echo response
	echoResponse := models.NewEchoResponse(echoRequest, "Request successfully echoed")

	// Convert response to JSON
	responseBody, err := echoResponse.ToJSON()
	if err != nil {
		h.logger.Error("Failed to marshal response", map[string]interface{}{
			"error": err.Error(),
		})
		return h.createErrorResponse(http.StatusInternalServerError, "Internal Server Error", "Failed to process response")
	}

	// Log the successful response with full response body
	h.logger.Info("Request successfully echoed", map[string]interface{}{
		"response_size": len(responseBody),
		"method":        request.HTTPMethod,
		"path":          request.Path,
		"response_body": responseBody,
		"message":       "Request processed successfully",
	})

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type, Authorization",
		},
		Body: responseBody,
	}, nil
}

// isMethodAllowed checks if the HTTP method is allowed
func (h *LambdaHandler) isMethodAllowed(method string) bool {
	allowedMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"OPTIONS": true, // For CORS preflight
	}
	return allowedMethods[method]
}

// parseRequest extracts request information from API Gateway proxy request
func (h *LambdaHandler) parseRequest(request *events.APIGatewayProxyRequest) *models.EchoRequest {
	// Convert headers to map[string]string
	headers := make(map[string]string)
	for key, value := range request.Headers {
		headers[key] = value
	}

	// Convert query parameters to map[string]string
	queryParams := make(map[string]string)
	for key, value := range request.QueryStringParameters {
		queryParams[key] = value
	}

	// Create the echo request
	return models.NewEchoRequest(
		request.HTTPMethod,
		request.Path,
		headers,
		queryParams,
		request.Body,
	)
}

// createErrorResponse creates a standardized error response
func (h *LambdaHandler) createErrorResponse(statusCode int, error, message string) (events.APIGatewayProxyResponse, error) {
	errorResponse := models.NewErrorResponse(error, message)
	
	responseBody, err := errorResponse.ToJSON()
	if err != nil {
		// Fallback to simple error response if JSON marshaling fails
		log.Printf("Failed to marshal error response: %v", err)
		responseBody = fmt.Sprintf(`{"error": "Internal Server Error", "message": "Failed to process error response", "timestamp": "%s"}`, errorResponse.Timestamp)
	}

	h.logger.Error("Error response generated", map[string]interface{}{
		"status_code":   statusCode,
		"error":         error,
		"message":       message,
		"response_body": responseBody,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type, Authorization",
		},
		Body: responseBody,
	}, nil
}