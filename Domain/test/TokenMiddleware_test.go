package test

import (
	"GOLANG/Domain/middleware"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TestJWTAuth_MissingToken tests JWT middleware without token
func TestJWTAuth_MissingToken(t *testing.T) {
	app := fiber.New()

	app.Use(middleware.JWTAuth())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("Protected route")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// TestJWTAuth_InvalidTokenFormat tests JWT middleware with invalid token format
func TestJWTAuth_InvalidTokenFormat(t *testing.T) {
	app := fiber.New()

	app.Use(middleware.JWTAuth())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("Protected route")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidToken")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// TestJWTAuth_BearerWithoutToken tests JWT middleware with Bearer but no token
func TestJWTAuth_BearerWithoutToken(t *testing.T) {
	app := fiber.New()

	app.Use(middleware.JWTAuth())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("Protected route")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// TestJWTAuth_InvalidToken tests JWT middleware with invalid JWT token
func TestJWTAuth_InvalidToken(t *testing.T) {
	app := fiber.New()

	app.Use(middleware.JWTAuth())
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("Protected route")
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// Note: Test for valid token requires actual JWT generation
// which will be tested in integration tests
