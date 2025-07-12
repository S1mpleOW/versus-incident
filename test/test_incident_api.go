package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Incident payload structure matching your curl request
type IncidentPayload struct {
	Action       string       `json:"action"`
	Data         Data         `json:"data"`
	Installation Installation `json:"installation"`
	Actor        Actor        `json:"actor"`
}

type Data struct {
	Issue Issue `json:"issue"`
}

type Issue struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Culprit   string   `json:"culprit"`
	ShortID   string   `json:"shortId"`
	Project   Project  `json:"project"`
	Metadata  Metadata `json:"metadata"`
	Status    string   `json:"status"`
	Level     string   `json:"level"`
	FirstSeen string   `json:"firstSeen"`
	LastSeen  string   `json:"lastSeen"`
	Count     int      `json:"count"`
	UserCount int      `json:"userCount"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Metadata struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Installation struct {
	UUID string `json:"uuid"`
}

type Actor struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Default to localhost, but allow override via command line
	baseURL := "http://localhost:3000"
	if len(os.Args) > 1 {
		baseURL = os.Args[1]
	}

	apiURL := baseURL + "/api/incidents"

	fmt.Printf("Testing Versus Incident API at: %s\n\n", apiURL)

	// Test different payloads
	fmt.Println("=== Test 1: Sentry-like Issue ===")
	testSentryIssue(apiURL)

	fmt.Println("\n=== Test 2: Critical Alert ===")
	testCriticalAlert(apiURL)

	fmt.Println("\n=== Test 3: Custom Alert ===")
	testCustomAlert(apiURL)
}

func testSentryIssue(apiURL string) {
	payload := IncidentPayload{
		Action: "created",
		Data: Data{
			Issue: Issue{
				ID:      "123456",
				Title:   "Example Issue",
				Culprit: "example_function in example_module",
				ShortID: "PROJECT-1",
				Project: Project{
					ID:   "1",
					Name: "Example Project",
					Slug: "example-project",
				},
				Metadata: Metadata{
					Type:  "ExampleError",
					Value: "This is an example error",
				},
				Status:    "unresolved",
				Level:     "error",
				FirstSeen: "2023-10-01T12:00:00Z",
				LastSeen:  "2023-10-01T12:05:00Z",
				Count:     5,
				UserCount: 3,
			},
		},
		Installation: Installation{
			UUID: "installation-uuid",
		},
		Actor: Actor{
			Type: "user",
			ID:   "789",
			Name: "John Doe",
		},
	}

	sendIncident(apiURL, payload, "Sentry Issue")
}

func testCriticalAlert(apiURL string) {
	payload := IncidentPayload{
		Action: "created",
		Data: Data{
			Issue: Issue{
				ID:      "987654",
				Title:   "Critical Database Connection Failure",
				Culprit: "database_connection in auth_service",
				ShortID: "CRITICAL-1",
				Project: Project{
					ID:   "2",
					Name: "Production Auth Service",
					Slug: "auth-service",
				},
				Metadata: Metadata{
					Type:  "DatabaseConnectionError",
					Value: "Connection timeout after 30 seconds",
				},
				Status:    "unresolved",
				Level:     "fatal",
				FirstSeen: time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
				LastSeen:  time.Now().Format(time.RFC3339),
				Count:     15,
				UserCount: 8,
			},
		},
		Installation: Installation{
			UUID: "prod-installation-uuid",
		},
		Actor: Actor{
			Type: "system",
			ID:   "monitoring",
			Name: "Production Monitoring",
		},
	}

	sendIncident(apiURL, payload, "Critical Alert")
}

func testCustomAlert(apiURL string) {
	// Test with a simpler custom payload
	customPayload := map[string]interface{}{
		"ServiceName": "test-service",
		"Message":     "Custom test alert for Google Chat integration",
		"Severity":    "HIGH",
		"Status":      "FIRING",
		"Timestamp":   time.Now().Format(time.RFC3339),
		"Environment": "development",
		"Source":      "test-script",
	}

	jsonData, _ := json.Marshal(customPayload)
	sendRawIncident(apiURL, jsonData, "Custom Alert")
}

func sendIncident(apiURL string, payload IncidentPayload, testName string) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("❌ Failed to marshal %s payload: %v\n", testName, err)
		return
	}

	sendRawIncident(apiURL, jsonData, testName)
}

func sendRawIncident(apiURL string, jsonData []byte, testName string) {
	fmt.Printf("Sending %s...\n", testName)
	fmt.Printf("Payload size: %d bytes\n", len(jsonData))
	fmt.Printf("Payload preview: %.200s...\n", string(jsonData))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "versus-incident-test/1.0")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	fmt.Printf("Sending request to: %s\n", req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Response Headers: %+v\n", resp.Header)
	fmt.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("✅ %s sent successfully!\n", testName)
	} else {
		fmt.Printf("❌ %s failed with status %d\n", testName, resp.StatusCode)

		// Troubleshooting tips
		if resp.StatusCode == 404 {
			fmt.Println("\n🔍 Troubleshooting 404 Error:")
			fmt.Println("1. Make sure Versus Incident is running on the specified port")
			fmt.Println("2. Check if the API endpoint path is correct (/api/incidents)")
			fmt.Println("3. Verify the server is accessible")
		} else if resp.StatusCode == 500 {
			fmt.Println("\n🔍 Troubleshooting 500 Error:")
			fmt.Println("1. Check Versus Incident logs for detailed error messages")
			fmt.Println("2. Verify Google Chat webhook URL is configured correctly")
			fmt.Println("3. Check if all required configuration is set")
		}
	}

	// Add a small delay between tests
	time.Sleep(2 * time.Second)
}
