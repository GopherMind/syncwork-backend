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
		return c.Status(400).JSON(fiber.Map{
			"error":   "INVALID_REQUEST_BODY",
			"message": "Invalid JSON format",
		})
	}

	if bodyUser.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "EMAIL_REQUIRED",
			"message": "Email is required",
		})
	}
	if bodyUser.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "PASSWORD_REQUIRED",
			"message": "Password is required",
		})
	}
	if len(bodyUser.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{
			"error":   "PASSWORD_TOO_SHORT",
			"message": "Password must be at least 6 characters",
		})
	}
	if bodyUser.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "NAME_REQUIRED",
			"message": "Name is required",
		})
	}
	if bodyUser.Role == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "ROLE_REQUIRED",
			"message": "Role is required",
		})
	}
	if bodyUser.Role != "freelancer" && bodyUser.Role != "client" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "INVALID_ROLE",
			"message": "Role must be 'freelancer' or 'client'",
		})
	}

	authResp, err := db.SB.Auth.Signup(types.SignupRequest{
		Email:    bodyUser.Email,
		Password: bodyUser.Password,
	})
	if err != nil {
		return c.Status(409).JSON(fiber.Map{
			"error":   "SIGNUP_FAILED",
			"message": "Email already exists or invalid credentials",
		})
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
		return c.Status(500).JSON(fiber.Map{
			"error":   "PROFILE_CREATION_FAILED",
			"message": "User created but profile setup failed",
		})
	}

	bodyUser.Id = authResp.User.ID.String()
	token, err := jwt.Createjwt(bodyUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "TOKEN_GENERATION_FAILED",
			"message": "Registration successful but token generation failed",
		})
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
