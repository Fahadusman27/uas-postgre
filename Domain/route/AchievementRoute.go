package route

import (
	"GOLANG/Domain/middleware"
	"GOLANG/Domain/service"

	"github.com/gofiber/fiber/v2"
)

// AchievementRoute - 5.4 Achievements
func AchievementRoute(API *fiber.App) {
	achievements := API.Group("/api/v1/achievements")

	// Semua endpoint butuh JWT authentication
	achievements.Use(middleware.JWTAuth())

	// GET /api/v1/achievements/advisee - View prestasi mahasiswa bimbingan (Dosen Wali)
	// Permission: verify_achievements
	// FR-006: View Prestasi Mahasiswa Bimbingan
	achievements.Get("/advisee",
		middleware.RequirePermission("verify_achievements"),
		service.GetAdviseeAchievementsService)

	// GET /api/v1/achievements - List all achievements (Admin)
	// Permission: read_achievements
	// FR-010: View All Achievements
	achievements.Get("/",
		middleware.RequirePermission("read_achievements"),
		service.GetAllAchievementsService)

	// GET /api/v1/achievements/:id - Detail achievement
	// Permission: read_achievements atau verify_achievements
	// achievements.Get("/:id",
	// 	middleware.RequireAnyPermission("read_achievements", "verify_achievements"),
	// 	service.GetAchievementDetailService)

	// POST /api/v1/achievements - Create achievement (Mahasiswa)
	// Permission: write_achievements
	achievements.Post("/",
		middleware.RequirePermission("write_achievements"),
		service.SubmitAchievementService)

	// PUT /api/v1/achievements/:id - Update achievement (Mahasiswa)
	// Permission: write_achievements
	// achievements.Put("/:id",
	// 	middleware.RequirePermission("write_achievements"),
	// 	service.UpdateAchievementService)

	// DELETE /api/v1/achievements/:id - Delete achievement (Mahasiswa)
	// Permission: write_achievements
	achievements.Delete("/:id",
		middleware.RequirePermission("write_achievements"),
		service.DeleteAchievementService)

	// POST /api/v1/achievements/:id/submit - Submit for verification
	// Permission: write_achievements
	achievements.Post("/:id/submit",
		middleware.RequirePermission("write_achievements"),
		service.SubmitForVerificationService)

	// POST /api/v1/achievements/:id/verify - Verify achievement (Dosen Wali)
	// Permission: verify_achievements
	// FR-007: Verify Prestasi
	achievements.Post("/:id/verify",
		middleware.RequirePermission("verify_achievements"),
		service.VerifyAchievementService)

	// POST /api/v1/achievements/:id/reject - Reject achievement (Dosen Wali)
	// Permission: verify_achievements
	// FR-008: Reject Prestasi
	achievements.Post("/:id/reject",
		middleware.RequirePermission("verify_achievements"),
		service.RejectAchievementService)

	// GET /api/v1/achievements/:id/history - Status history
	// Permission: read_achievements atau verify_achievements
	// achievements.Get("/:id/history",
	// 	middleware.RequireAnyPermission("read_achievements", "verify_achievements"),
	// 	service.GetAchievementHistoryService)

	// POST /api/v1/achievements/:id/attachments - Upload files
	// Permission: write_achievements
	// achievements.Post("/:id/attachments",
	// 	middleware.RequirePermission("write_achievements"),
	// 	service.UploadAttachmentsService)
}
