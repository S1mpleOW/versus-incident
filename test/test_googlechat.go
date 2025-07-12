package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Simple test message for Google Chat
type SimpleGoogleChatMessage struct {
	Text string `json:"text"`
}

// Cards v2 format test message
type CardsV2Message struct {
	CardsV2 []Card `json:"cardsV2"`
}

type Card struct {
	Card CardContent `json:"card"`
}

type CardContent struct {
	Header   *Header   `json:"header,omitempty"`
	Sections []Section `json:"sections"`
}

type Header struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
}

type Section struct {
	Widgets []Widget `json:"widgets"`
}

type Widget struct {
	TextParagraph *TextParagraph `json:"textParagraph,omitempty"`
}

type TextParagraph struct {
	Text string `json:"text"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_googlechat.go <webhook_url>")
		fmt.Println("Example: go run test_googlechat.go 'https://chat.googleapis.com/v1/spaces/.../messages?key=...'")
		os.Exit(1)
	}

	webhookURL := os.Args[1]

	fmt.Printf("Testing Google Chat webhook URL: %s\n\n", webhookURL)

	// Test 1: Simple text message
	fmt.Println("=== Test 1: Simple Text Message ===")
	testSimpleMessage(webhookURL)

	fmt.Println("\n=== Test 2: Cards v2 Message ===")
	testCardsV2Message(webhookURL)
}

func testSimpleMessage(webhookURL string) {
	msg := SimpleGoogleChatMessage{
		Text: "🔥 Test alert from Versus Incident - Simple Format",
	}

	jsonData, _ := json.Marshal(msg)
	sendTestMessage(webhookURL, jsonData, "Simple Text")
}

func testCardsV2Message(webhookURL string) {
	msg := CardsV2Message{
		CardsV2: []Card{{
			Card: CardContent{
				Header: &Header{
					Title:    "🔥 Test Alert from Versus Incident",
					Subtitle: "Cards v2 Format Test",
				},
				Sections: []Section{{
					Widgets: []Widget{{
						TextParagraph: &TextParagraph{
							Text: "<b>Status:</b> FIRING<br><b>Severity:</b> CRITICAL<br><b>Description:</b> This is a test alert to verify Google Chat integration is working correctly.",
						},
					}},
				}},
			},
		}},
	}

	jsonData, _ := json.Marshal(msg)
	sendTestMessage(webhookURL, jsonData, "Cards v2")
}

func sendTestMessage(webhookURL string, jsonData []byte, messageType string) {
	fmt.Printf("Sending %s message...\n", messageType)
	fmt.Printf("Payload: %s\n", string(jsonData))

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Failed to create request: %v\n", err)
		return
	}

	// Set headers that Google Chat expects
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", "versus-incident-test/1.0")

	fmt.Printf("Request URL: %s\n", req.URL.String())
	fmt.Printf("Request Headers: %+v\n", req.Header)

	client := &http.Client{}
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
		fmt.Printf("✅ %s message sent successfully!\n", messageType)
	} else {
		fmt.Printf("❌ %s message failed with status %d\n", messageType, resp.StatusCode)

		// Common troubleshooting tips
		if resp.StatusCode == 404 {
			fmt.Println("\n🔍 Troubleshooting 404 Error:")
			fmt.Println("1. Verify your webhook URL is complete and correct")
			fmt.Println("2. Make sure the Google Chat app/bot is added to the space")
			fmt.Println("3. Check if the webhook URL has the correct format:")
			fmt.Println("   https://chat.googleapis.com/v1/spaces/{space}/messages?key={key}&token={token}")
			fmt.Println("4. Ensure the space ID and key/token are valid")
		}
	}
}
