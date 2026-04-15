package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/GopherMind/syncwork-backend/utils/jwt"
)

func AuthMiddleware (c *fiber.Ctx)error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Токен отсутствует",
		})
	}

	userId, err := jwt.CheckJwt(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "токен",
		})
	}

	c.Locals("user_id", userId)

	return  c.Next()
}