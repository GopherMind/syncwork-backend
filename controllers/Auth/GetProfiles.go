package auth

import (
	"log"
	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/models"
	"github.com/gofiber/fiber/v2"
)


func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "UNAUTHORIZED",
			"message": "User ID not found in context",
		})
	}

	var profile models.Profile 
	var userTasks []models.Task
	var proposals []models.Proposal
	_, err := db.SB.From("profiles").
		Select("*", "", false).
		Eq("id", userID.(string)).
		Single().
		ExecuteTo(&profile) 
	if err != nil {
		log.Printf("Error fetching profile: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to fetch profile",
		})
	}

	_, err = db.SB.From("tasks").
		Select("*", "", false).
		Eq("client_id", userID.(string)).
		ExecuteTo(&userTasks)
	if err != nil {
		log.Printf("Error fetching user tasks: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to fetch user tasks",
		})
	}
	_, err = db.SB.From("propals").
		Select("*", "", false).
		Eq("user_id", userID.(string)).
		ExecuteTo(&proposals)
	if err != nil {
		log.Printf("Error fetching proposals: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to fetch proposals",
		})
	}
	return c.JSON(fiber.Map{
		"profile": profile,
		"tasks":   userTasks,
		"proposals": proposals,
	})
}