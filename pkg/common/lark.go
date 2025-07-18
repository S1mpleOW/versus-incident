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

type LarkProvider struct {
	webhookURL   string
	templatePath string
	client       *http.Client
}

func NewLarkProvider(cfg config.LarkConfig, proxyConfig config.ProxyConfig) *LarkProvider {
	client := utils.CreateHTTPClient(proxyConfig, cfg.UseProxy)
	return &LarkProvider{
		webhookURL:   cfg.WebhookURL,
		templatePath: cfg.TemplatePath,
		client:       client,
	}
}

func (l *LarkProvider) SendAlert(i *m.Incident) error {
	funcMaps := utils.GetTemplateFuncMaps()

	tmpl, err := template.New(filepath.Base(l.templatePath)).Funcs(funcMaps).ParseFiles(l.templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var message bytes.Buffer
	if err := tmpl.Execute(&message, i.Content); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create interactive card message
	larkMsg := utils.CreateLarkMessage(message.String(), i.Resolved)

	jsonData, err := json.Marshal(larkMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", l.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("lark API returned non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
