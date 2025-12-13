package route

import (
	"GOLANG/Domain/middleware"

	"github.com/gofiber/fiber/v2"
)

// AuthRoute - 5.1 Authentication (Tanpa Handler Eksplisit)
func AuthRoute(API *fiber.App) {
	auth := API.Group("/api/v1/auth")

	// POST /api/v1/auth/login - Public route
	auth.Post("/login", middleware.CallService("AuthService", "Login"))

	// POST /api/v1/auth/refresh - Public route (TODO: implement refresh token)
	auth.Post("/refresh", middleware.CallService("AuthService", "Refresh"))

	// POST /api/v1/auth/logout - Protected route
	auth.Post("/logout",
		middleware.JWTAuth(),
		middleware.CallService("AuthService", "Logout"),
	)

	// GET /api/v1/auth/profile - Protected route
	auth.Get("/profile",
		middleware.JWTAuth(),
		middleware.CallService("AuthService", "GetProfile"),
	)
}