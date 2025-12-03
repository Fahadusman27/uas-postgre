package route

import (
	"GOLANG/Domain/middleware"

	"github.com/gofiber/fiber/v2"
)

// UserRoute - 5.2 Users (Admin)
func UserRoute(app *fiber.App) {
	API := app.Group("/api/v1")

	// Semua endpoint users butuh JWT authentication
	API.Use(middleware.JWTAuth())

	// GET /api/v1/users - List all users
	// Permission: read_users
	// users.Get("/",
	// 	middleware.RequirePermission("read_users"),
	// 	service.GetAllUsersService)

	// GET /api/v1/users/:id - Get user by ID
	// Permission: read_users
	// users.Get("/:id",
	// 	middleware.RequirePermission("read_users"),
	// 	service.GetUserByIDService)

	// POST /api/v1/users - Create new user
	// Permission: write_users
	// users.Post("/",
	// 	middleware.RequirePermission("write_users"),
	// 	service.CreateUserService)

	// PUT /api/v1/users/:id - Update user
	// Permission: write_users
	// users.Put("/:id",
	// 	middleware.RequirePermission("write_users"),
	// 	service.UpdateUserService)

	// DELETE /api/v1/users/:id - Delete user
	// Permission: write_users
	// users.Delete("/:id",
	// 	middleware.RequirePermission("write_users"),
	// 	service.DeleteUserService)

	// PUT /api/v1/users/:id/role - Update user role
	// Permission: write_users
	// users.Put("/:id/role",
	// 	middleware.RequirePermission("write_users"),
	// 	service.UpdateUserRoleService)
}
