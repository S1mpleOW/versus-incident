package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/VersusControl/versus-incident/pkg/common"
	"github.com/VersusControl/versus-incident/pkg/config"
	"github.com/VersusControl/versus-incident/pkg/core"
	"github.com/VersusControl/versus-incident/pkg/dto"
	m "github.com/VersusControl/versus-incident/pkg/models"
)

// CreateNotification sends notifications based on the provided payload
func CreateNotification(source string, payload *map[string]interface{}, overwriteValues ...map[string]string) error {
	var cfg *config.Config

	if len(overwriteValues) > 0 {
		cfg = config.GetConfigWithParamsOverwrite(&overwriteValues[0])
	} else {
		cfg = config.GetConfig()
	}

	// Handle Bitbucket build payload specifically
	if isValidBitbucketPayload(*payload) {
		return handleBitbucketNotification(cfg, payload)
	}
	// Handle generic notifications
	return handleGenericNotification(cfg, payload)
}

// isValidBitbucketPayload checks if the payload is a valid Bitbucket build payload
func isValidBitbucketPayload(payload map[string]interface{}) bool {
	requiredFields := []string{"type", "key", "state", "name", "url", "description"}

	for _, field := range requiredFields {
		if _, exists := payload[field]; !exists {
			return false
		}
	}

	// Check if type is "build"
	if payloadType, ok := payload["type"].(string); !ok || payloadType != "build" {
		return false
	}

	return true
}

// handleBitbucketNotification processes Bitbucket build notifications
func handleBitbucketNotification(cfg *config.Config, payload *map[string]interface{}) error {
	// Parse payload into Bitbucket DTO
	payloadBytes, err := json.Marshal(*payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	var bitbucketPayload dto.BitbucketBuildPayload
	if err := json.Unmarshal(payloadBytes, &bitbucketPayload); err != nil {
		return fmt.Errorf("failed to unmarshal Bitbucket payload: %v", err)
	}

	// Create notification message based on build state
	utc7Location := time.FixedZone("UTC+7", 7*60*60) // UTC+7 timezone
	formattedTime := time.Now().In(utc7Location).Format("2006-01-02 15:04:05")

	notificationData := map[string]interface{}{
		"title":       fmt.Sprintf("Bitbucket Build %s", capitalizeState(bitbucketPayload.State)),
		"description": bitbucketPayload.Description,
		"build_key":   bitbucketPayload.Key,
		"build_name":  bitbucketPayload.Name,
		"build_state": bitbucketPayload.State,
		"build_url":   bitbucketPayload.URL,
		"severity":    determineSeverity(bitbucketPayload.State),
		"timestamp":   formattedTime,
	}

	// Create incident for the notification
	incident := m.NewIncident("bitbucket", &notificationData, isResolvedState(bitbucketPayload.State))

	// Send notifications using the alert factory pattern
	return sendNotificationsViaAlert(cfg, incident)
}

// handleGenericNotification processes generic notifications
func handleGenericNotification(cfg *config.Config, payload *map[string]interface{}) error {
	// Create incident for generic notification
	incident := m.NewIncident("generic", payload, false)

	// Send notifications using the alert factory pattern
	return sendNotificationsViaAlert(cfg, incident)
}

// sendNotificationsViaAlert sends notifications using the existing alert factory pattern
func sendNotificationsViaAlert(cfg *config.Config, incident *m.Incident) error {
	// Initialize providers using the existing factory
	factory := common.NewAlertProviderFactory(cfg)
	providers, err := factory.CreateProviders()
	if err != nil {
		return fmt.Errorf("failed to create providers: %v", err)
	}

	// Create alert with providers
	alert := core.NewAlert(providers...)

	// Send the alert
	if err := alert.SendAlert(incident); err != nil {
		return fmt.Errorf("failed to send notifications: %v", err)
	}

	return nil
}

// Utility functions
func capitalizeState(state string) string {
	if len(state) == 0 {
		return state
	}
	return string(state[0]-32) + state[1:]
}

func determineSeverity(state string) string {
	switch state {
	case "FAILED":
		return "critical"
	case "STOPPED":
		return "error"
	case "SUCCESSFUL":
		return "resolved"
	case "INPROGRESS":
		return "warning"
	default:
		return "unknown"
	}
}

func isResolvedState(state string) bool {
	return state == "SUCCESSFUL"
}
