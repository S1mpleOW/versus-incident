package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/VersusControl/versus-incident/pkg/dto"
)

// TestBitbucketPayload tests the Bitbucket payload structure
func TestBitbucketPayload(t *testing.T) {
	// Sample payload based on user requirements
	payload := dto.BitbucketBuildPayload{
		Type:        "build",
		Key:         "BUILD-123",
		State:       "SUCCESSFUL",
		Name:        "Model Alias changed",
		URL:         "https://bitbucket.org/repo/builds/123",
		Description: "Build completed successfully",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	// Verify all required fields are present
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		t.Fatalf("Failed to unmarshal payload: %v", err)
	}

	requiredFields := []string{"type", "key", "state", "name", "url", "description"}
	for _, field := range requiredFields {
		if _, exists := data[field]; !exists {
			t.Errorf("Missing required field: %s", field)
		}
	}

	// Verify type is "build"
	if data["type"] != "build" {
		t.Errorf("Expected type 'build', got '%v'", data["type"])
	}
}

// Example of how to use the endpoint
func ExampleBitbucketRequest() {
	// This is the payload structure that Bitbucket pipeline should send
	payload := map[string]interface{}{
		"type":        "build",
		"key":         "$BITBUCKET_BUILD_NUMBER",
		"state":       "$STATUS",
		"name":        "Model Alias changed",
		"url":         "$BUILD_URL",
		"description": "$DESC",
	}

	jsonData, _ := json.Marshal(payload)

	// POST to the Bitbucket notification endpoint
	resp, err := http.Post(
		"http://localhost:8080/api/notifications/bitbucket",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		// Handle error
		return
	}
	defer resp.Body.Close()

	// Response should be 201 Created with success message
}
