package handler

import (
	"context"
	"encoding/json"
	"log"

	"echo-api/internal/models"
	"echo-api/pkg/logger"
)

// NonProxyRequest represents the request structure for non-proxy integration
type NonProxyRequest struct {
	HTTPMethod               string            `json:"httpMethod"`
	Path                     string            `json:"path"`
	Headers                  map[string]string `json:"headers"`
	QueryStringParameters    map[string]string `json:"queryStringParameters"`
	Body                     string            `json:"body"`
}

// NonProxyHandler handles AWS Lambda non-proxy requests
type NonProxyHandler struct {
	logger *logger.Logger
}

// NewNonProxyHandler creates a new non-proxy Lambda handler instance
func NewNonProxyHandler() *NonProxyHandler {
	return &NonProxyHandler{
		logger: logger.New(),
	}
}

// HandleRequest processes the incoming non-proxy request
func (h *NonProxyHandler) HandleRequest(ctx context.Context, request NonProxyRequest) (map[string]interface{}, error) {
	h.logger.Info("Processing non-proxy request", map[string]interface{}{
		"method":     request.HTTPMethod,
		"path":       request.Path,
		"headers":    request.Headers,
		"query":      request.QueryStringParameters,
		"body":       request.Body,
	})

	// Check if method is allowed (GET or POST)
	if !h.isMethodAllowed(request.HTTPMethod) {
		h.logger.Warn("Method not allowed", map[string]interface{}{
			"method": request.HTTPMethod,
		})
		return h.createErrorResponse(405, "Method Not Allowed", "Only GET and POST methods are supported")
	}

	// Parse the request
	echoRequest := h.parseRequest(&request)

	// Create echo response
	echoResponse := models.NewEchoResponse(echoRequest, "Request successfully echoed")

	// Convert response to map
	responseMap := map[string]interface{}{
		"request":     echoResponse.Request,
		"message":     echoResponse.Message,
		"processedAt": echoResponse.ProcessedAt,
	}

	// Convert response to JSON for logging
	responseJSON, err := json.Marshal(responseMap)
	if err != nil {
		h.logger.Error("Failed to marshal response", map[string]interface{}{
			"error": err.Error(),
		})
		return h.createErrorResponse(500, "Internal Server Error", "Failed to process response")
	}

	// Log the successful response with full response body
	h.logger.Info("Request processed successfully", map[string]interface{}{
		"response_size": len(responseJSON),
		"method":        request.HTTPMethod,
		"path":          request.Path,
		"response_body": string(responseJSON),
	})

	return responseMap, nil
}

// isMethodAllowed checks if the HTTP method is allowed
func (h *NonProxyHandler) isMethodAllowed(method string) bool {
	allowedMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"OPTIONS": true, // For CORS preflight
	}
	return allowedMethods[method]
}

// parseRequest extracts request information from non-proxy request
func (h *NonProxyHandler) parseRequest(request *NonProxyRequest) *models.EchoRequest {
	// Create the echo request
	return models.NewEchoRequest(
		request.HTTPMethod,
		request.Path,
		request.Headers,
		request.QueryStringParameters,
		request.Body,
	)
}

// createErrorResponse creates a standardized error response
func (h *NonProxyHandler) createErrorResponse(statusCode int, errorType, message string) (map[string]interface{}, error) {
	errorResponse := models.NewErrorResponse(errorType, message)
	
	responseMap := map[string]interface{}{
		"error":     errorResponse.Error,
		"message":   errorResponse.Message,
		"timestamp": errorResponse.Timestamp,
	}

	responseJSON, err := json.Marshal(responseMap)
	if err != nil {
		// Fallback to simple error response if JSON marshaling fails
		log.Printf("Failed to marshal error response: %v", err)
		responseMap = map[string]interface{}{
			"error":     "Internal Server Error",
			"message":   "Failed to process error response",
			"timestamp": errorResponse.Timestamp,
		}
		responseJSON, _ = json.Marshal(responseMap)
	}

	h.logger.Error("Error response generated", map[string]interface{}{
		"status_code":   statusCode,
		"error":         errorType,
		"message":       message,
		"response_body": string(responseJSON),
	})

	return responseMap, nil
}