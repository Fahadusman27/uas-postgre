package route

import (
	"GOLANG/Domain/middleware"
	"GOLANG/Domain/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App) {
	auth := app.Group("/auth")

	// Public routes
	auth.Post("/login", service.LoginService)

	// Protected routes
	auth.Post("/logout", middleware.JWTAuth(), service.LogoutService)
}
