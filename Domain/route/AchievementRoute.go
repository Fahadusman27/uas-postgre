package route

import (
	"GOLANG/Domain/middleware"
	"GOLANG/Domain/service"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoute(app *fiber.App) {
	achievements := app.Group("/achievements")

	// Semua endpoint butuh JWT authentication
	achievements.Use(middleware.JWTAuth())

	// POST /achievements - Submit prestasi (FR-003)
	// Actor: Mahasiswa
	// Permission: write_achievements
	achievements.Post("/",
		middleware.RequirePermission("write_achievements"),
		service.SubmitAchievementService)

	// Endpoint lainnya (untuk future implementation)
	// GET /achievements - List achievements
	// achievements.Get("/",
	// 	middleware.RequireAnyPermission("read_achievements", "verify_achievements"),
	// 	service.GetAchievementsService)

	// GET /achievements/:id - Detail achievement
	// achievements.Get("/:id",
	// 	middleware.RequireAnyPermission("read_achievements", "verify_achievements"),
	// 	service.GetAchievementDetailService)

	// PUT /achievements/:id - Update achievement (mahasiswa only)
	// achievements.Put("/:id",
	// 	middleware.RequirePermission("write_achievements"),
	// 	service.UpdateAchievementService)

	// DELETE /achievements/:id - Delete achievement (mahasiswa only)
	// achievements.Delete("/:id",
	// 	middleware.RequirePermission("write_achievements"),
	// 	service.DeleteAchievementService)

	// PUT /achievements/:id/submit - Submit untuk verifikasi
	// achievements.Put("/:id/submit",
	// 	middleware.RequirePermission("write_achievements"),
	// 	service.SubmitForVerificationService)

	// PUT /achievements/:id/verify - Verify achievement (dosen/admin only)
	// achievements.Put("/:id/verify",
	// 	middleware.RequirePermission("verify_achievements"),
	// 	service.VerifyAchievementService)

	// PUT /achievements/:id/reject - Reject achievement (dosen/admin only)
	// achievements.Put("/:id/reject",
	// 	middleware.RequirePermission("verify_achievements"),
	// 	service.RejectAchievementService)
}
