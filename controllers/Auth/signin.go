package auth

import (
	"encoding/json"
	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/models"
	"github.com/GopherMind/syncwork-backend/utils/jwt"
	"github.com/gofiber/fiber/v2"
)


func Signin(c *fiber.Ctx) error {
	var bodyUser models.UserAuth

	if err := c.BodyParser(&bodyUser); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "INVALID_REQUEST_BODY",
            "message": "Invalid JSON format",
        })
    }

	if bodyUser.Email == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "EMAIL_REQUIRED",
            "message": "Email is required",
        })
    }
    if bodyUser.Password == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "PASSWORD_REQUIRED",
            "message": "Password is required",
        })
    }
    if len(bodyUser.Password) < 6 {
        return c.Status(400).JSON(fiber.Map{
            "error": "PASSWORD_TOO_SHORT",
            "message": "Password must be at least 6 characters",
        })
    }

	resp, err := db.SB.Auth.SignInWithEmailPassword(bodyUser.Email, bodyUser.Password)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "USER_NOT_FOUND",
            "message": "users is not found",
        })
    }

    var result []map[string]interface{}
    data, _, err := db.SB.From("profiles").Select("name", "", false).Eq("id", resp.User.ID.String()).Execute()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "PROFILE_FETCH_FAILED",
            "message": "Failed to fetch user profile",
        })
    }

    if err := json.Unmarshal(data, &result); err != nil || len(result) == 0 {
        return c.Status(500).JSON(fiber.Map{
            "error": "PROFILE_PARSE_FAILED",
            "message": "Failed to parse user profile",
        })
    }

    bodyUser.Id = resp.User.ID.String()
    bodyUser.Name = result[0]["name"].(string)
    token, err := jwt.Createjwt(bodyUser)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "TOKEN_GENERATION_FAILED",
            "message": "Failed to generate token",
        })
    }

    return c.Status(200).JSON(fiber.Map{"token": token})
}
