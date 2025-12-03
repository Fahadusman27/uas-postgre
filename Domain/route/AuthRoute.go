package route

import (
	"GOLANG/Domain/middleware"
	"GOLANG/Domain/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App) {
	auth := app.Group("/auth")

	// Public routes - tidak perlu authentication
	auth.Post("/login", service.LoginService)

	// Protected routes - perlu JWT authentication
	// Logout hanya perlu JWT valid, tidak perlu permission khusus
	auth.Post("/logout", middleware.JWTAuth(), service.LogoutService)
}

// Contoh route dengan RBAC (Permission-based)
// Uncomment dan sesuaikan dengan service yang ada
/*
func UserRoute(app *fiber.App) {
	users := app.Group("/users")

	// Semua endpoint users butuh JWT authentication
	users.Use(middleware.JWTAuth())

	// GET /users - butuh permission "read_users"
	users.Get("/",
		middleware.RequirePermission("read_users"),
		service.GetAllUsers)

	// POST /users - butuh permission "write_users"
	users.Post("/",
		middleware.RequirePermission("write_users"),
		service.CreateUser)

	// PUT /users/:id - butuh permission "write_users"
	users.Put("/:id",
		middleware.RequirePermission("write_users"),
		service.UpdateUser)

	// DELETE /users/:id - butuh permission "write_users"
	users.Delete("/:id",
		middleware.RequirePermission("write_users"),
		service.DeleteUser)
}

func AchievementRoute(app *fiber.App) {
	achievements := app.Group("/achievements")

	// Semua endpoint achievements butuh JWT authentication
	achievements.Use(middleware.JWTAuth())

	// GET /achievements - butuh salah satu: read_achievements atau verify_achievements
	achievements.Get("/",
		middleware.RequireAnyPermission("read_achievements", "verify_achievements"),
		service.GetAchievements)

	// POST /achievements - butuh permission "write_achievements"
	achievements.Post("/",
		middleware.RequirePermission("write_achievements"),
		service.CreateAchievement)

	// PUT /achievements/:id - butuh permission "write_achievements"
	achievements.Put("/:id",
		middleware.RequirePermission("write_achievements"),
		service.UpdateAchievement)

	// DELETE /achievements/:id - butuh permission "write_achievements"
	achievements.Delete("/:id",
		middleware.RequirePermission("write_achievements"),
		service.DeleteAchievement)

	// PUT /achievements/:id/verify - butuh permission "verify_achievements"
	achievements.Put("/:id/verify",
		middleware.RequirePermission("verify_achievements"),
		service.VerifyAchievement)

	// PUT /achievements/:id/reject - butuh permission "verify_achievements"
	achievements.Put("/:id/reject",
		middleware.RequirePermission("verify_achievements"),
		service.RejectAchievement)
}
*/
