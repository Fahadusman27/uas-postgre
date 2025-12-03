package service

import (
	"GOLANG/Domain/config"
	model "GOLANG/Domain/model/Postgresql"
	. "GOLANG/Domain/repository"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginService(c *fiber.Ctx) error {
	var body model.Login

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Cari user berdasarkan email atau username
	var user *model.Users
	var err error

	if body.Email != "" {
		user, err = GetUserByEmail(body.Email)
	} else if body.Username != "" {
		user, err = GetUserByUsername(body.Username)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email atau username harus diisi",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email atau username salah",
		})
	}

	// Cek status aktif user
	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Akun dinonaktifkan",
		})
	}

	// Validasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Password salah",
		})
	}

	// Ambil permissions berdasarkan role
	permissions, err := GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data permissions",
		})
	}

	// Generate JWT token
	jwtSecret := []byte(config.GetJWTSecret())
	expiryTime := time.Now().Add(config.GetJWTExpiry())

	claims := jwt.MapClaims{
		"id":          user.ID.String(),
		"username":    user.Username,
		"role_id":     user.RoleID.String(),
		"permissions": permissions,
		"exp":         expiryTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat token session",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": fiber.Map{
			"id":        user.ID,
			"username":  user.Username,
			"full_name": user.FullName,
			"email":     user.Email,
			"role_id":   user.RoleID,
		},
	})
}

func LogoutService(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token tidak ditemukan",
		})
	}

	// Extract token dari header
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Parse token untuk mendapatkan expiry time
	jwtSecret := []byte(config.GetJWTSecret())
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token invalid",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membaca claims",
		})
	}

	// Ambil expiry time dari claims
	exp := int64(claims["exp"].(float64))
	expiresAt := time.Unix(exp, 0)

	// Simpan token ke blacklist
	if time.Now().Before(expiresAt) {
		err := AddTokenToBlacklist(tokenString, expiresAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal logout, silakan coba lagi",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout berhasil",
	})
}
