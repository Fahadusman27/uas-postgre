package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireRole middleware untuk cek role_id
// Digunakan ketika endpoint memerlukan role tertentu (role-based access)
func RequireRole(allowedRoleIDs ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil role_id dari context (sudah di-set oleh JWTAuth middleware)
		roleID := c.Locals("role_id")
		if roleID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Role tidak ditemukan",
			})
		}

		roleIDStr, ok := roleID.(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Format role tidak valid",
			})
		}

		// Check apakah role_id ada di allowed list
		for _, allowedID := range allowedRoleIDs {
			if roleIDStr == allowedID {
				// Allow request - user memiliki role yang sesuai
				return c.Next()
			}
		}

		// Deny request - user tidak memiliki role yang sesuai
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Akses ditolak",
			"message": "Role Anda tidak diizinkan mengakses endpoint ini",
		})
	}
}
