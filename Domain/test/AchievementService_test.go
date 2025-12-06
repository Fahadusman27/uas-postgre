package test

import (
	"GOLANG/Domain/service"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TestSubmitAchievementService_InvalidRequestBody tests submit achievement with invalid body
func TestSubmitAchievementService_InvalidRequestBody(t *testing.T) {
	app := fiber.New()

	// Mock JWT middleware - set user ID in context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("id", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	})

	app.Post("/achievements", service.SubmitAchievementService)

	// Invalid JSON
	req := httptest.NewRequest("POST", "/achievements", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestSubmitAchievementService_MissingTitle tests submit achievement without title
func TestSubmitAchievementService_MissingTitle(t *testing.T) {
	app := fiber.New()

	// Mock JWT middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("id", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	})

	app.Post("/achievements", service.SubmitAchievementService)

	achievementData := map[string]interface{}{
		"achievementType": "competition",
		"description":     "Test description",
	}
	body, _ := json.Marshal(achievementData)

	req := httptest.NewRequest("POST", "/achievements", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestSubmitAchievementService_InvalidAchievementType tests invalid achievement type
func TestSubmitAchievementService_InvalidAchievementType(t *testing.T) {
	app := fiber.New()

	// Mock JWT middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("id", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	})

	app.Post("/achievements", service.SubmitAchievementService)

	achievementData := map[string]interface{}{
		"achievementType": "invalid_type",
		"title":           "Test Achievement",
		"description":     "Test description",
	}
	body, _ := json.Marshal(achievementData)

	req := httptest.NewRequest("POST", "/achievements", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestSubmitForVerificationService_InvalidAchievementID tests invalid achievement ID
func TestSubmitForVerificationService_InvalidAchievementID(t *testing.T) {
	app := fiber.New()

	// Mock JWT middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("id", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	})

	app.Post("/achievements/:id/submit", service.SubmitForVerificationService)

	req := httptest.NewRequest("POST", "/achievements/invalid-id/submit", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestRejectAchievementService_MissingRejectionNote tests reject without note
func TestRejectAchievementService_MissingRejectionNote(t *testing.T) {
	app := fiber.New()

	// Mock JWT middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("id", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	})

	app.Post("/achievements/:id/reject", service.RejectAchievementService)

	// Valid MongoDB ObjectID format
	achievementID := "507f1f77bcf86cd799439011"

	rejectData := map[string]string{
		"rejection_note": "",
	}
	body, _ := json.Marshal(rejectData)

	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/reject", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
