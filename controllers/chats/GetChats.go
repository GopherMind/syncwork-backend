package chats

import (
	"log"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/gofiber/fiber/v2"
)

type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
}

type Chat struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	TaskID string `json:"task_id"`
	Task   Task   `json:"tasks"`
}
type ChatUser struct {
	ChatID string `json:"chat_id"`
	Chat   Chat   `json:"chats"`
}

func GetChats(c *fiber.Ctx) error {
	userID := c.Locals("user_id")

	if userID == nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "UNAUTHORIZED",
			"message": "User ID not found in context",
		})
	}

	var chatUsers []ChatUser

	_, err := db.SB.From("chat_users").
		Select("chat_id, chats(id, status, task_id, tasks(id, title, description))", "", false).
		Eq("user_id", userID.(string)).
		ExecuteTo(&chatUsers)

	if err != nil {
		log.Println("Error fetching chats:", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"chats": chatUsers,
	})
}
