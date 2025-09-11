package logger

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"
)

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		output: log.New(&buf, "", 0),
	}

	message := "Test info message"
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	logger.Info(message, data)

	output := buf.String()
	if output == "" {
		t.Error("Expected log output, got empty string")
	}

	// Parse the JSON log entry
	var entry LogEntry
	err := json.Unmarshal([]byte(output), &entry)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	if entry.Level != INFO {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}
	if entry.Message != message {
		t.Errorf("Expected message %s, got %s", message, entry.Message)
	}
	if entry.Data["key1"] != "value1" {
		t.Errorf("Expected data key1=value1, got %v", entry.Data["key1"])
	}
	if entry.Timestamp == "" {
		t.Error("Expected timestamp to be set")
	}
}

func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		output: log.New(&buf, "", 0),
	}

	message := "Test warning message"
	data := map[string]interface{}{
		"warning_code": 1001,
	}

	logger.Warn(message, data)

	output := buf.String()
	var entry LogEntry
	err := json.Unmarshal([]byte(output), &entry)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	if entry.Level != WARN {
		t.Errorf("Expected level WARN, got %s", entry.Level)
	}
	if entry.Message != message {
		t.Errorf("Expected message %s, got %s", message, entry.Message)
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		output: log.New(&buf, "", 0),
	}

	message := "Test error message"
	data := map[string]interface{}{
		"error_code": "E500",
		"details":    "Something went wrong",
	}

	logger.Error(message, data)

	output := buf.String()
	var entry LogEntry
	err := json.Unmarshal([]byte(output), &entry)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	if entry.Level != ERROR {
		t.Errorf("Expected level ERROR, got %s", entry.Level)
	}
	if entry.Message != message {
		t.Errorf("Expected message %s, got %s", message, entry.Message)
	}
	if entry.Data["error_code"] != "E500" {
		t.Errorf("Expected error_code E500, got %v", entry.Data["error_code"])
	}
}

func TestLogger_WithNilData(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		output: log.New(&buf, "", 0),
	}

	logger.Info("Test message", nil)

	output := buf.String()
	var entry LogEntry
	err := json.Unmarshal([]byte(output), &entry)
	if err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	if entry.Data != nil {
		t.Errorf("Expected data to be nil, got %v", entry.Data)
	}
}