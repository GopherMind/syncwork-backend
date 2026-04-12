package auth

import (

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/models"
	"github.com/GopherMind/syncwork-backend/utils/hashPassword"
	"github.com/GopherMind/syncwork-backend/utils/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/gotrue-go/types"
)

func Login(c *fiber.Ctx) error {
	var bodyUser models.UserAuth
	if err := c.BodyParser(&bodyUser); err != nil {

		return c.Status(400).JSON(fiber.Map{"error": "Auth failed: " + err.Error()})

	}
	if bodyUser.Email == "" || bodyUser.Password == "" || bodyUser.Name == "" || bodyUser.Role == "" {
		return c.Status(400).JSON(fiber.Map{"error": "empty fields"})

	}
	hashedPassword, err := hashpassword.HashPassword(bodyUser.Password)
	authResp, err := db.SB.Auth.Signup(types.SignupRequest{
		Email:    bodyUser.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Auth failed: " + err.Error()})
	}

	newProfile := map[string]interface{}{
		"id":          authResp.User.ID,     // ID из Auth
		"name":        bodyUser.Name,        // Колонка name
		"description": bodyUser.Description, // Колонка description
		"url":         bodyUser.Url,         // Колонка url
		"role":        bodyUser.Role,        // Колонка role
	}

	_, _, err = db.SB.From("profiles").Insert(newProfile, false, "", "minimal", "").Execute()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Auth failed: " + err.Error()})
	}
	token, err := jwt.Createjwt(bodyUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "jwt create error: " + err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{
		"token": token,
	})
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