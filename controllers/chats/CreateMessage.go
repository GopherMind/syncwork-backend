package chats

import (
	"log"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/gofiber/fiber/v2"
)

func CreateMessage(c *fiber.Ctx) error {

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "UNAUTHORIZED",
			"message": "User ID not found in context",
		})
	}

	chatID := c.Params("id")
	if chatID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "BAD_REQUEST",
			"message": "Chat ID is required",
		})
	}
	type Request struct {
		ChatID   string `json:"chat_id"`
		Message  string `json:"message"`
		SenderID string `json:"sender_id"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error":   "BAD_REQUEST",
			"message": "Invalid request body",
		})
	}

	req.ChatID = chatID
	req.SenderID = userID.(string)

	_, _, err := db.SB.From("messages").Insert(&req, false, "", "", "").Execute()
	if err != nil {
		log.Printf("Database error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "DATABASE_ERROR",
			"message": "Failed to create message",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Message created successfully",
	})
}
