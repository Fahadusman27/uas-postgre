package route

import (
	"GOLANG/Domain/middleware"

	"github.com/gofiber/fiber/v2"
)

func StudentRoute(API *fiber.App) {
	// Students endpoints
	students := API.Group("/api/v1/students")
	students.Use(middleware.JWTAuth())

	// GET /api/v1/students - List all students
	students.Get("/",
		middleware.RequirePermission("read_students"),
		middleware.CallService("StudentService", "GetStudents"))

	// GET /api/v1/students/:id - Get student detail
	students.Get("/:id",
		middleware.RequirePermission("read_students"),
		middleware.CallService("StudentService", "GetStudentDetail"))

	// GET /api/v1/students/:id/achievements - Get student achievements
	students.Get("/:id/achievements",
		middleware.RequireAnyPermission("read_achievements", "verify_achievements"),
		middleware.CallService("StudentService", "GetStudentAchievements"))

	// PUT /api/v1/students/:id/advisor - Set student advisor
	students.Put("/:id/advisor",
		middleware.RequirePermission("manage_students"),
		middleware.CallService("StudentService", "SetStudentAdvisor"))

	// Lecturers endpoints
	lecturers := API.Group("/api/v1/lecturers")
	lecturers.Use(middleware.JWTAuth())

	// GET /api/v1/lecturers - List all lecturers
	lecturers.Get("/",
		middleware.RequirePermission("read_lecturers"),
		middleware.CallService("LecturerService", "GetLecturers"))

	// GET /api/v1/lecturers/:id/advisees - Get lecturer advisees
	lecturers.Get("/:id/advisees",
		middleware.RequirePermission("read_lecturers"),
		middleware.CallService("LecturerService", "GetLecturerAdvisees"))
}
