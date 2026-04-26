package propasals

import (
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
		Tasks  struct {
			ClientID string `json:"client_id"`
		} `json:"tasks"`
		Chats struct {
			ID string `json:"id"`
		} `json:"chats"`
	}

	_, err := db.SB.From("proposals").
		Select("id, task_id, tasks(client_id), chats", "", false).
		Eq("id", proposalID).
		ExecuteTo(&result)
	if err != nil || len(result) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error":   "NOT_FOUND",
			"message": "Proposal not found",
		})
	}

	if result[0].Tasks.ClientID != userID.(string) {
		return c.Status(403).JSON(fiber.Map{
			"error":   "FORBIDDEN",
			"message": "Only task owner can accept proposals",
		})
	}

	_, _, err = db.SB.From("proposals"). 
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

	_, _, err = db.SB.From("chat_user").Insert(map[string]interface{}{
		"chat_id": result[0].Chats.ID,
		"user_id": result[0].Tasks.ClientID,
	}, false, "", "", "").Execute()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to update task status",
		})
	}

	return c.JSON(fiber.Map{"message": "Proposal accepted and task status updated"})


}

