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

// TestLoginService_InvalidRequestBody tests login with invalid request body
func TestLoginService_InvalidRequestBody(t *testing.T) {
	app := fiber.New()
	app.Post("/login", service.LoginService)

	// Invalid JSON
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestLoginService_MissingFields tests login with missing required fields
func TestLoginService_MissingFields(t *testing.T) {
	app := fiber.New()
	app.Post("/login", service.LoginService)

	// Missing password
	loginData := map[string]string{
		"email": "test@example.com",
	}
	body, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestLoginService_InvalidEmail tests login with invalid email format
func TestLoginService_InvalidEmail(t *testing.T) {
	app := fiber.New()
	app.Post("/login", service.LoginService)

	loginData := map[string]string{
		"email":    "invalid-email",
		"password": "password123",
	}
	body, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// Note: Tests for successful login require database mocking
// which will be implemented in integration tests
