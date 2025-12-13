package route

import (
	"GOLANG/Domain/middleware"

	"github.com/gofiber/fiber/v2"
)

// ReportRoute - 5.8 Reports & Analytics (Tanpa Handler Eksplisit)
func ReportRoute(API *fiber.App) {
	reports := API.Group("/api/v1/reports")
	reports.Use(middleware.JWTAuth())

	// GET /api/v1/reports/statistics - General statistics
	// Permission based on role: students see own, advisors see advisees, admins see all
	reports.Get("/statistics",
		middleware.CallService("ReportService", "GetStatistics"))

	// GET /api/v1/reports/student/:id - Specific student statistics
	// Permission: read_achievements or verify_achievements
	reports.Get("/student/:id",
		middleware.RequireAnyPermission("read_achievements", "verify_achievements"),
		middleware.CallService("ReportService", "GetStudentStatistics"))
}
