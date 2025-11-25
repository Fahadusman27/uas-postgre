package service

import (
	"POSTGRE/Domain/model"

	"github.com/gofiber/fiber/v2"
)


func AuthService(c *fiber.Ctx) error {
	var body struct {
		model.Login
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid Body"})
	}

	user, err := 
}