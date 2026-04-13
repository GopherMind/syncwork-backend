package auth

import (

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/models"
	"github.com/GopherMind/syncwork-backend/utils/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/gotrue-go/types"
)

func Register(c *fiber.Ctx) error {
    var bodyUser models.UserAuth
    if err := c.BodyParser(&bodyUser); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Parse failed: " + err.Error()})
    }

    if bodyUser.Email == "" || bodyUser.Password == "" || bodyUser.Name == "" || bodyUser.Role == "" {
        return c.Status(400).JSON(fiber.Map{"error": "empty fields"})
    }
	if bodyUser.Role != "freelancer" && bodyUser.Role != "client" {
		return c.Status(400).JSON(fiber.Map{"error": "unknown role"})
	}
    authResp, err := db.SB.Auth.Signup(types.SignupRequest{
        Email:    bodyUser.Email,
        Password: bodyUser.Password,
    })
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Signup failed: " + err.Error()})
    }

    newProfile := map[string]interface{}{
        "id":          authResp.User.ID,
        "name":        bodyUser.Name,
        "description": bodyUser.Description,
        "url":         bodyUser.Url,
        "role":        bodyUser.Role,
    }

    _, _, err = db.SB.From("profiles").Insert(newProfile, false, "", "minimal", "").Execute()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Profile insert failed: " + err.Error()})
    }

    token, err := jwt.Createjwt(bodyUser)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "JWT error: " + err.Error()})
    }

    return c.Status(201).JSON(fiber.Map{"token": token})
}

/*
пример запроса:
{
	"email": "test@test.com",
	"password": "123456"
	"name": "test",
	"role": "test",
}
*/