package middleware

import (
	"github.com/GopherMind/syncwork-backend/utils/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Токен отсутствует",
		})
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	println("Token after cleanup:", token)
	claims, err := jwt.CheckJwt(token)

	if err != nil {
		println("JWT Check Error:", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "wrong token",
			"details": err.Error(),
		})
	}
	println("Claims ID:", claims.Id)

	c.Locals("user_id", claims.Id)

	return c.Next()
}
