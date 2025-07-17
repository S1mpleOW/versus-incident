package controllers

import (
	"fmt"

	"github.com/VersusControl/versus-incident/pkg/config"
	"github.com/VersusControl/versus-incident/pkg/services"

	"github.com/gofiber/fiber/v2"
)

func CreateNotification(c *fiber.Ctx) error {
	cfg := config.GetConfig()

	if cfg.Alert.DebugBody {
		rawBody := c.Body()

		// Log the raw request body for debugging purposes
		fmt.Println("Raw Request Body:", string(rawBody))
	}

	body := &map[string]interface{}{}

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var err error

	// If query parameters exist, get the value to overwrite the default configuration
	if len(c.Queries()) > 0 {
		overwriteValue := c.Queries()
		err = services.CreateNotification("", body, overwriteValue)
	} else {
		err = services.CreateNotification("", body)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "Notification sent successfully"})
}

// CreateBitbucketNotification handles Bitbucket pipeline notifications specifically
func CreateBitbucketNotification(c *fiber.Ctx) error {
	cfg := config.GetConfig()

	if cfg.Alert.DebugBody {
		rawBody := c.Body()
		fmt.Println("Raw Bitbucket Request Body:", string(rawBody))
	}

	body := &map[string]interface{}{}

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Bitbucket payload"})
	}

	// Validate required Bitbucket fields
	requiredFields := []string{"type", "key", "state", "name", "url", "description"}
	for _, field := range requiredFields {
		if _, exists := (*body)[field]; !exists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Missing required field: %s", field),
			})
		}
	}

	var err error

	// If query parameters exist, get the value to overwrite the default configuration
	if len(c.Queries()) > 0 {
		overwriteValue := c.Queries()
		err = services.CreateNotification("bitbucket", body, overwriteValue)
	} else {
		err = services.CreateNotification("bitbucket", body)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "Bitbucket notification sent successfully",
		"build":  (*body)["key"],
		"state":  (*body)["state"],
	})
}
