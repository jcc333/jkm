package log

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	// Clean up before test
	os.Remove("test.jsonl")

	// Create a logger that logs to a test file
	logger, err := New(true)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Change the file path to test file
	logger.file, err = os.OpenFile("test.jsonl", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}

	// Log a test message
	err = logger.Info("test message")
	if err != nil {
		t.Fatalf("Failed to log message: %v", err)
	}

	// Close the logger
	err = logger.Close()
	if err != nil {
		t.Fatalf("Failed to close logger: %v", err)
	}

	// Read the log file
	data, err := os.ReadFile("test.jsonl")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Parse the JSON
	var entry logEntry
	err = json.Unmarshal(data, &entry)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check the entry
	if entry.Message != "test message" {
		t.Errorf("Expected message to be 'test message', got '%s'", entry.Message)
	}

	if entry.Level != LevelInfo {
		t.Errorf("Expected level to be 'info', got '%s'", entry.Level)
	}

	// Clean up after test
	os.Remove("test.jsonl")
}
