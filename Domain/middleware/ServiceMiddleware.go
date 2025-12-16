package middleware

import (
	"GOLANG/Domain/service"

	"github.com/gofiber/fiber/v2"
)

// ServiceMiddleware untuk auto-call service methods tanpa handler eksplisit
func CallService(serviceName, methodName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		switch serviceName {
		case "AuthService":
			return callAuthService(c, methodName)
		case "UserService":
			return callUserService(c, methodName)
		case "AchievementService":
			return callAchievementService(c, methodName)
		case "StudentService":
			return callStudentService(c, methodName)
		case "LecturerService":
			return callLecturerService(c, methodName)
		case "ReportService":
			return callReportService(c, methodName)
		default:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Service not found: " + serviceName,
			})
		}
	}
}

// Auth Service Calls
func callAuthService(c *fiber.Ctx, methodName string) error {
	switch methodName {
	case "Login":
		return service.LoginService(c)
	case "Logout":
		return service.LogoutService(c)
	case "Refresh":
		// TODO: Implement refresh token service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Refresh token not implemented yet",
		})
	case "GetProfile":
		// TODO: Implement get profile service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Get profile not implemented yet",
		})
	default:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Method not found: " + methodName,
		})
	}
}

// User Service Calls
func callUserService(c *fiber.Ctx, methodName string) error {
	switch methodName {
	case "CreateUser":
		return service.CreateUserService(c)
	case "GetUsers":
		return service.GetUsersService(c)
	case "GetUserDetail":
		return service.GetUserDetailService(c)
	case "UpdateUser":
		return service.UpdateUserService(c)
	case "DeleteUser":
		return service.DeleteUserService(c)
	case "AssignPassword":
		return service.UpdateUserPasswordService(c)
	case "SetStudentProfile":
		return service.SetStudentProfileService(c)
	case "SetLecturerProfile":
		return service.SetLecturerProfileService(c)
	default:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Method not found: " + methodName,
		})
	}
}

// Achievement Service Calls
func callAchievementService(c *fiber.Ctx, methodName string) error {
	switch methodName {
	case "SubmitAchievement":
		return service.SubmitAchievementService(c)
	case "SubmitForVerification":
		return service.SubmitForVerificationService(c)
	case "DeleteAchievement":
		return service.DeleteAchievementService(c)
	case "GetAdviseeAchievements":
		return service.GetAdviseeAchievementsService(c)
	case "VerifyAchievement":
		return service.VerifyAchievementService(c)
	case "RejectAchievement":
		return service.RejectAchievementService(c)
	case "GetAllAchievements":
		return service.GetAllAchievementsService(c)
	case "GetMyAchievementStats":
		return service.GetMyAchievementStatsService(c)
	case "GetAdviseeAchievementStats":
		return service.GetAdviseeAchievementStatsService(c)
	case "GetAllAchievementStats":
		return service.GetAllAchievementStatsService(c)
	case "GetAchievementDetail":
		// TODO: Implement get achievement detail service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Get achievement detail not implemented yet",
		})
	case "UpdateAchievement":
		// TODO: Implement update achievement service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Update achievement not implemented yet",
		})
	case "GetAchievementHistory":
		// TODO: Implement get achievement history service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Get achievement history not implemented yet",
		})
	case "UploadAttachments":
		// TODO: Implement upload attachments service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Upload attachments not implemented yet",
		})
	default:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Method not found: " + methodName,
		})
	}
}

// Student Service Calls
func callStudentService(c *fiber.Ctx, methodName string) error {
	switch methodName {
	case "GetStudents":
		// TODO: Implement get students service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Get students not implemented yet",
		})
	case "GetStudentDetail":
		// TODO: Implement get student detail service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Get student detail not implemented yet",
		})
	case "GetStudentAchievements":
		// TODO: Implement get student achievements service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Get student achievements not implemented yet",
		})
	case "SetStudentAdvisor":
		// TODO: Implement set student advisor service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Set student advisor not implemented yet",
		})
	default:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Method not found: " + methodName,
		})
	}
}

// Lecturer Service Calls
func callLecturerService(c *fiber.Ctx, methodName string) error {
	switch methodName {
	case "GetLecturers":
		// TODO: Implement get lecturers service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Get lecturers not implemented yet",
		})
	case "GetLecturerAdvisees":
		// TODO: Implement get lecturer advisees service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Get lecturer advisees not implemented yet",
		})
	default:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Method not found: " + methodName,
		})
	}
}

// Report Service Calls
func callReportService(c *fiber.Ctx, methodName string) error {
	switch methodName {
	case "GetStatistics":
		// TODO: Implement get statistics service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Get statistics not implemented yet",
		})
	case "GetStudentStatistics":
		// TODO: Implement get student statistics service
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Get student statistics not implemented yet",
		})
	default:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Method not found: " + methodName,
		})
	}
}
