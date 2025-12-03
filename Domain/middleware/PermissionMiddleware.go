package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequirePermission middleware untuk cek single permission
// Digunakan ketika endpoint memerlukan 1 permission spesifik
func RequirePermission(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil permissions dari context (sudah di-set oleh JWTAuth middleware)
		perms := c.Locals("permissions")
		if perms == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Permission tidak ditemukan",
			})
		}

		// 2. Convert ke []interface{} (dari JWT claims)
		permissions, ok := perms.([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Format permission tidak valid",
			})
		}

		// 3. Check apakah user memiliki permission yang diperlukan
		for _, perm := range permissions {
			if permStr, ok := perm.(string); ok && permStr == requiredPermission {
				// Allow request - user memiliki permission
				return c.Next()
			}
		}

		// 4. Deny request - user tidak memiliki permission
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":               "Akses ditolak",
			"required_permission": requiredPermission,
			"message":             "Anda tidak memiliki permission: " + requiredPermission,
		})
	}
}

// RequireAnyPermission middleware untuk cek salah satu dari beberapa permission
// Digunakan ketika endpoint bisa diakses dengan salah satu permission (OR logic)
func RequireAnyPermission(requiredPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		perms := c.Locals("permissions")
		if perms == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Permission tidak ditemukan",
			})
		}

		permissions, ok := perms.([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Format permission tidak valid",
			})
		}

		// Check apakah user memiliki salah satu permission yang diperlukan
		for _, perm := range permissions {
			permStr, ok := perm.(string)
			if !ok {
				continue
			}
			for _, required := range requiredPermissions {
				if permStr == required {
					// Allow request - user memiliki salah satu permission
					return c.Next()
				}
			}
		}

		// Deny request - user tidak memiliki satupun permission yang diperlukan
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":                "Akses ditolak",
			"required_permissions": requiredPermissions,
			"message":              "Anda tidak memiliki salah satu permission yang diperlukan",
		})
	}
}

// RequireAllPermissions middleware untuk cek semua permission harus ada
// Digunakan ketika endpoint memerlukan multiple permissions (AND logic)
func RequireAllPermissions(requiredPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		perms := c.Locals("permissions")
		if perms == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Permission tidak ditemukan",
			})
		}

		permissions, ok := perms.([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Format permission tidak valid",
			})
		}

		// Convert user permissions ke map untuk lookup cepat
		userPerms := make(map[string]bool)
		for _, perm := range permissions {
			if permStr, ok := perm.(string); ok {
				userPerms[permStr] = true
			}
		}

		// Check apakah user memiliki semua permission yang diperlukan
		var missingPermissions []string
		for _, required := range requiredPermissions {
			if !userPerms[required] {
				missingPermissions = append(missingPermissions, required)
			}
		}

		if len(missingPermissions) > 0 {
			// Deny request - user tidak memiliki semua permission
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":                "Akses ditolak",
				"missing_permissions":  missingPermissions,
				"required_permissions": requiredPermissions,
				"message":              "Anda tidak memiliki semua permission yang diperlukan",
			})
		}

		// Allow request - user memiliki semua permission
		return c.Next()
	}
}
