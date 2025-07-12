package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/VersusControl/versus-incident/pkg/config"
	m "github.com/VersusControl/versus-incident/pkg/models"
	"github.com/VersusControl/versus-incident/pkg/utils"
)

type GoogleChatProvider struct {
	webhookURL   string
	templatePath string
	client       *http.Client
}

func NewGoogleChatProvider(cfg config.GoogleChatConfig, proxyConfig config.ProxyConfig) *GoogleChatProvider {
	client := utils.CreateHTTPClient(proxyConfig, cfg.UseProxy)
	return &GoogleChatProvider{
		webhookURL:   cfg.WebhookURL,
		templatePath: cfg.TemplatePath,
		client:       client,
	}
}

func (g *GoogleChatProvider) SendAlert(i *m.Incident) error {
	funcMaps := utils.GetTemplateFuncMaps()

	tmpl, err := template.New(filepath.Base(g.templatePath)).Funcs(funcMaps).ParseFiles(g.templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var message bytes.Buffer
	if err := tmpl.Execute(&message, i.Content); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// The template should already generate valid Google Chat JSON format
	// So we need to parse it to ensure it's valid JSON and then send it
	var googleChatMsg interface{}
	if err := json.Unmarshal(message.Bytes(), &googleChatMsg); err != nil {
		return fmt.Errorf("failed to parse template output as JSON: %w", err)
	}

	// Re-marshal to ensure clean JSON
	jsonData, err := json.Marshal(googleChatMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Enhanced debugging
	fmt.Printf("=== Google Chat Debug Info ===\n")
	fmt.Printf("Webhook URL: %s\n", g.webhookURL)
	fmt.Printf("Payload size: %d bytes\n", len(jsonData))
	fmt.Printf("Payload preview (first 500 chars): %.500s\n", string(jsonData))

	req, err := http.NewRequest("POST", g.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers for Google Chat
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", "versus-incident/1.0")

	fmt.Printf("Request headers: %+v\n", req.Header)
	fmt.Printf("Request method: %s\n", req.Method)
	fmt.Printf("Request URL: %s\n", req.URL.String())

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response status: %d\n", resp.StatusCode)
	fmt.Printf("Response headers: %+v\n", resp.Header)
	fmt.Printf("Response body: %s\n", string(body))
	fmt.Printf("=== End Debug Info ===\n")

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("google Chat API returned non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
