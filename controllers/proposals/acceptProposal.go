package propasals

import (
	"fmt"
	"log"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/gofiber/fiber/v2"
)

func AcceptProposal(c *fiber.Ctx) error {
	proposalID := c.Params("id")
	userID := c.Locals("user_id")

	if proposalID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "BAD_REQUEST",
			"message": "Proposal ID is required",
		})
	}

	if userID == nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "UNAUTHORIZED",
			"message": "User ID not found in context",
		})
	}

	var result []struct {
		ID     string `json:"id"`
		TaskID string `json:"task_id"`
		UserID string `json:"user_id"`
		Status string `json:"status"`
	}
	log.Printf("Accepting proposal ID: %s by user: %s", proposalID, userID)
	_, err := db.SB.From("propals").
		Select("id, task_id, user_id, status", "", false).
		Eq("id", proposalID).
		Eq("status", "pending").
		ExecuteTo(&result)

	if err != nil {
		log.Printf("Database error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "DATABASE_ERROR",
			"message": fmt.Sprintf("Database error: %v", err),
		})
	}

	if len(result) == 0 {
		log.Printf("Proposal not found or not pending: %s", proposalID)
		return c.Status(404).JSON(fiber.Map{
			"error":   "NOT_FOUND",
			"message": "Proposal not found or already processed",
		})
	}

	log.Printf("Found proposal: %+v", result[0])

	var tasks []struct {
		ClientID string `json:"client_id"`
	}
	_, err = db.SB.From("tasks").
		Select("client_id", "", false).
		Eq("id", result[0].TaskID).
		ExecuteTo(&tasks)

	if err != nil || len(tasks) == 0 {
		log.Printf("Task not found: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error":   "NOT_FOUND",
			"message": "Task not found",
		})
	}

	if tasks[0].ClientID != userID.(string) {
		return c.Status(403).JSON(fiber.Map{
			"error":   "FORBIDDEN",
			"message": "Only task owner can accept proposals",
		})
	}

	_, _, err = db.SB.From("propals").
		Update(map[string]interface{}{"status": "accepted"}, "", "").
		Eq("id", proposalID).
		Execute()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to update proposal",
		})
	}
	_, _, err = db.SB.From("tasks").
		Update(map[string]interface{}{"status": "closed"}, "", "").
		Eq("id", result[0].TaskID).
		Execute()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to update task status",
		})
	}

	var chats []struct {
		ID string `json:"id"`
	}
	_, err = db.SB.From("chats").
		Select("id", "", false).
		Eq("task_id", result[0].TaskID).
		ExecuteTo(&chats)

	if err != nil || len(chats) == 0 {
		log.Printf("Chat lookup error: %v, chats found: %d", err, len(chats))
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to find chat for task",
		})
	}

	_, _, err = db.SB.From("chat_user").Insert(map[string]interface{}{
		"chat_id": chats[0].ID,
		"user_id": result[0].UserID,
	}, false, "", "", "").Execute()
	if err != nil {
		log.Printf("Failed to add user to chat: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to add user to chat",
		})
	}

	_, _, err = db.SB.From("propals").Update(map[string]interface{}{"status": "denied"}, "", "").
		Eq("task_id", result[0].TaskID).
		Neq("id", proposalID).
		Execute()
	if err != nil {
		log.Printf("Failed to deny other proposals: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to update other proposals status",
		})
	}
	return c.JSON(fiber.Map{"message": "Proposal accepted and task status updated"})

}


