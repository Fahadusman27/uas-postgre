package route

import (
	"GOLANG/Domain/middleware"

	"github.com/gofiber/fiber/v2"
)

// AchievementRoute - 5.4 Achievements (Tanpa Handler Eksplisit)
func AchievementRoute(API *fiber.App) {
	achievements := API.Group("/api/v1/achievements")

	// Semua endpoint butuh JWT authentication
	achievements.Use(middleware.JWTAuth())

	// GET /api/v1/achievements/stats/my - Statistics prestasi sendiri (Mahasiswa)
	// Permission: write_achievements
	// FR-011: Achievement Statistics (Own)
	achievements.Get("/stats/my", middleware.RequirePermission("write_achievements"),
		middleware.CallService("AchievementService", "GetMyAchievementStats"))

	// GET /api/v1/achievements/stats/advisee - Statistics prestasi mahasiswa bimbingan (Dosen Wali)
	// Permission: verify_achievements
	// FR-011: Achievement Statistics (Advisee)
	achievements.Get("/stats/advisee", middleware.RequirePermission("verify_achievements"),
		middleware.CallService("AchievementService", "GetAdviseeAchievementStats"))

	// GET /api/v1/achievements/stats/all - Statistics semua prestasi (Admin)
	// Permission: read_achievements
	// FR-011: Achievement Statistics (All)
	achievements.Get("/stats/all", middleware.RequirePermission("read_achievements"),
		middleware.CallService("AchievementService", "GetAllAchievementStats"))

	// GET /api/v1/achievements/advisee - View prestasi mahasiswa bimbingan (Dosen Wali)
	// Permission: verify_achievements
	// FR-006: View Prestasi Mahasiswa Bimbingan
	achievements.Get("/advisee", middleware.RequirePermission("verify_achievements"),
		middleware.CallService("AchievementService", "GetAdviseeAchievements"))

	// GET /api/v1/achievements - List all achievements (Admin)
	// Permission: read_achievements
	// FR-010: View All Achievements
	achievements.Get("/", middleware.RequirePermission("read_achievements"),
		middleware.CallService("AchievementService", "GetAllAchievements"))

	// GET /api/v1/achievements/:id - Detail achievement
	// Permission: read_achievements atau verify_achievements
	achievements.Get("/:id",
		middleware.RequireAnyPermission("read_achievements", "verify_achievements"),
		middleware.CallService("AchievementService", "GetAchievementDetail"))

	// POST /api/v1/achievements - Create achievement (Mahasiswa)
	// Permission: write_achievements
	achievements.Post("/", middleware.RequirePermission("write_achievements"),
		middleware.CallService("AchievementService", "SubmitAchievement"))

	// PUT /api/v1/achievements/:id - Update achievement (Mahasiswa)
	// Permission: write_achievements
	achievements.Put("/:id",
		middleware.RequirePermission("write_achievements"),
		middleware.CallService("AchievementService", "UpdateAchievement"))

	// DELETE /api/v1/achievements/:id - Delete achievement (Mahasiswa)
	// Permission: write_achievements
	achievements.Delete("/:id", middleware.RequirePermission("write_achievements"),
		middleware.CallService("AchievementService", "DeleteAchievement"))

	// POST /api/v1/achievements/:id/submit - Submit for verification
	// Permission: write_achievements
	achievements.Post("/:id/submit", middleware.RequirePermission("write_achievements"),
		middleware.CallService("AchievementService", "SubmitForVerification"))

	// POST /api/v1/achievements/:id/verify - Verify achievement (Dosen Wali)
	// Permission: verify_achievements
	// FR-007: Verify Prestasi
	achievements.Post("/:id/verify", middleware.RequirePermission("verify_achievements"),
		middleware.CallService("AchievementService", "VerifyAchievement"))

	// POST /api/v1/achievements/:id/reject - Reject achievement (Dosen Wali)
	// Permission: verify_achievements
	// FR-008: Reject Prestasi
	achievements.Post("/:id/reject", middleware.RequirePermission("verify_achievements"),
		middleware.CallService("AchievementService", "RejectAchievement"))

	// GET /api/v1/achievements/:id/history - Status history
	// Permission: read_achievements atau verify_achievements
	achievements.Get("/:id/history",
		middleware.RequireAnyPermission("read_achievements", "verify_achievements"),
		middleware.CallService("AchievementService", "GetAchievementHistory"))

	// POST /api/v1/achievements/:id/attachments - Upload files
	// Permission: write_achievements
	achievements.Post("/:id/attachments",
		middleware.RequirePermission("write_achievements"),
		middleware.CallService("AchievementService", "UploadAttachments"))
}
