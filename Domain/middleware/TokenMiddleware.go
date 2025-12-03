package middleware

import (
	"strings"

	"GOLANG/Domain/config"
	"GOLANG/Domain/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header required"})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header format must be Bearer {token}"})
		}

		tokenString := parts[1]

		// Cek apakah token ada di blacklist
		isBlacklisted, err := repository.IsTokenBlacklisted(tokenString)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memvalidasi token"})
		}
		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token telah di-logout"})
		}

		// Parse dan validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTSecret()), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Ambil data dari claims
		var userID string
		var roleID string
		var username string
		var permissions []interface{}

		if v, exists := claims["id"].(string); exists {
			userID = v
		}

		if r, exists := claims["role_id"].(string); exists {
			roleID = r
		}

		if u, exists := claims["username"].(string); exists {
			username = u
		}

		if p, exists := claims["permissions"].([]interface{}); exists {
			permissions = p
		}

		// Simpan ke context
		c.Locals("id", userID)
		c.Locals("role_id", roleID)
		c.Locals("username", username)
		c.Locals("permissions", permissions)

		return c.Next()
	}
}
