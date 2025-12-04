package route

import (
	"GOLANG/Domain/middleware"
	"GOLANG/Domain/service"

	"github.com/gofiber/fiber/v2"
)

// UserRoute - FR-009: Manage Users
func UserRoute(API *fiber.App) {
	users := API.Group("/api/v1/users")

	// Semua endpoint butuh JWT authentication dan permission manage_users
	users.Use(middleware.JWTAuth())
	users.Use(middleware.RequirePermission("manage_users"))

	// POST /api/v1/users - Create user
	users.Post("/",
		service.CreateUserService)

	// GET /api/v1/users - List users dengan pagination
	users.Get("/",
		service.GetUsersService)

	// GET /api/v1/users/:id - Get user detail
	users.Get("/:id",
		service.GetUserDetailService)

	// PUT /api/v1/users/:id - Update user
	users.Put("/:id",
		service.UpdateUserService)

	// DELETE /api/v1/users/:id - Delete user
	users.Delete("/:id",
		service.DeleteUserService)

	// PUT /api/v1/users/:id/role - Assign role
	users.Put("/:id/role",
		service.AssignRoleService)

	// POST /api/v1/users/:id/student - Set student profile
	users.Post("/:id/student",
		service.SetStudentProfileService)

	// POST /api/v1/users/:id/lecturer - Set lecturer profile
	users.Post("/:id/lecturer",
		service.SetLecturerProfileService)
}
