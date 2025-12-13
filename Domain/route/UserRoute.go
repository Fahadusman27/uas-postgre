package route

import (
	"GOLANG/Domain/middleware"

	"github.com/gofiber/fiber/v2"
)

// UserRoute - FR-009: Manage Users (Tanpa Handler Eksplisit)
func UserRoute(API *fiber.App) {
	users := API.Group("/api/v1/users")

	// Semua endpoint butuh JWT authentication dan permission manage_users
	users.Use(middleware.JWTAuth())
	users.Use(middleware.RequirePermission("manage_users"))

	// POST /api/v1/users - Create user
	users.Post("/",
		middleware.CallService("UserService", "CreateUser"))

	// GET /api/v1/users - List users dengan pagination
	users.Get("/",
		middleware.CallService("UserService", "GetUsers"))

	// GET /api/v1/users/:id - Get user detail
	users.Get("/:id",
		middleware.CallService("UserService", "GetUserDetail"))

	// PUT /api/v1/users/:id - Update user
	users.Put("/:id",
		middleware.CallService("UserService", "UpdateUser"))

	// DELETE /api/v1/users/:id - Delete user
	users.Delete("/:id",
		middleware.CallService("UserService", "DeleteUser"))

	// PUT /api/v1/users/:id/role - Assign role
	users.Put("/:id/role",
		middleware.CallService("UserService", "AssignRole"))

	// POST /api/v1/users/:id/student - Set student profile
	users.Post("/:id/student",
		middleware.CallService("UserService", "SetStudentProfile"))

	// POST /api/v1/users/:id/lecturer - Set lecturer profile
	users.Post("/:id/lecturer",
		middleware.CallService("UserService", "SetLecturerProfile"))
}
