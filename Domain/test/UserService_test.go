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

// TestCreateUserService_InvalidRequestBody tests create user with invalid body
func TestCreateUserService_InvalidRequestBody(t *testing.T) {
	app := fiber.New()
	app.Post("/users", service.CreateUserService)

	req := httptest.NewRequest("POST", "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestCreateUserService_MissingRequiredFields tests create user with missing fields
func TestCreateUserService_MissingRequiredFields(t *testing.T) {
	testCases := []struct {
		name string
		data map[string]string
	}{
		{
			name: "Missing username",
			data: map[string]string{
				"full_name": "Test User",
				"email":     "test@example.com",
				"password":  "password123",
				"role_id":   "550e8400-e29b-41d4-a716-446655440000",
			},
		},
		{
			name: "Missing email",
			data: map[string]string{
				"username":  "testuser",
				"full_name": "Test User",
				"password":  "password123",
				"role_id":   "550e8400-e29b-41d4-a716-446655440000",
			},
		},
		{
			name: "Missing password",
			data: map[string]string{
				"username":  "testuser",
				"full_name": "Test User",
				"email":     "test@example.com",
				"role_id":   "550e8400-e29b-41d4-a716-446655440000",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/users", service.CreateUserService)

			body, _ := json.Marshal(tc.data)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})
	}
}

// TestCreateUserService_InvalidRoleID tests create user with invalid role ID
func TestCreateUserService_InvalidRoleID(t *testing.T) {
	app := fiber.New()
	app.Post("/users", service.CreateUserService)

	userData := map[string]string{
		"username":  "testuser",
		"full_name": "Test User",
		"email":     "test@example.com",
		"password":  "password123",
		"role_id":   "invalid-uuid",
	}
	body, _ := json.Marshal(userData)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestSetStudentProfileService_MissingRequiredFields tests set student profile with missing fields
func TestSetStudentProfileService_MissingRequiredFields(t *testing.T) {
	app := fiber.New()
	app.Post("/users/:id/student", service.SetStudentProfileService)

	// Missing program_study
	profileData := map[string]string{
		"student_id":    "NIM123",
		"academic_year": "2023/2024",
	}
	body, _ := json.Marshal(profileData)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	req := httptest.NewRequest("POST", "/users/"+userID+"/student", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestSetLecturerProfileService_MissingRequiredFields tests set lecturer profile with missing fields
func TestSetLecturerProfileService_MissingRequiredFields(t *testing.T) {
	app := fiber.New()
	app.Post("/users/:id/lecturer", service.SetLecturerProfileService)

	// Missing department
	profileData := map[string]string{
		"lecturer_id": "NIDN123",
	}
	body, _ := json.Marshal(profileData)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	req := httptest.NewRequest("POST", "/users/"+userID+"/lecturer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
