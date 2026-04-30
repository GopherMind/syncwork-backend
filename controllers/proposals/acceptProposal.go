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
		ExecuteTo(&result)

	if err != nil {
		log.Printf("Database error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "DATABASE_ERROR",
			"message": fmt.Sprintf("Database error: %v", err),
		})
	}

	log.Printf("Raw query result (len=%d): %+v", len(result), result)

	if len(result) == 0 {
		log.Printf("Proposal not found: %s", proposalID)
		return c.Status(404).JSON(fiber.Map{
			"error":   "NOT_FOUND",
			"message": "Proposal not found",
		})
	}

	if result[0].Status != "pending" {
		log.Printf("Proposal %s has wrong status: '%s'", proposalID, result[0].Status)
		return c.Status(409).JSON(fiber.Map{
			"error":   "CONFLICT",
			"message": fmt.Sprintf("Proposal is already '%s'", result[0].Status),
		})
	}

	proposal := result[0]
	log.Printf("Found proposal: %+v", proposal)

	var tasks []struct {
		ClientID string `json:"client_id"`
	}

	_, err = db.SB.From("tasks").
		Select("client_id", "", false).
		Eq("id", proposal.TaskID).
		ExecuteTo(&tasks)

	if err != nil {
		log.Printf("Database error fetching task: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "DATABASE_ERROR",
			"message": "Failed to fetch task",
		})
	}

	if len(tasks) == 0 {
		log.Printf("Task not found: %s", proposal.TaskID)
		return c.Status(404).JSON(fiber.Map{
			"error":   "NOT_FOUND",
			"message": "Task not found",
		})
	}

	if tasks[0].ClientID != userID.(string) {
		log.Printf("Forbidden: user %s is not owner of task %s (owner: %s)",
			userID, proposal.TaskID, tasks[0].ClientID)
		return c.Status(403).JSON(fiber.Map{
			"error":   "FORBIDDEN",
			"message": "Only task owner can accept proposals",
		})
	}

	_, _, err = db.SB.From("propals").
		Update(map[string]interface{}{"status": "accepted"}, "", "").
		Eq("id", proposalID).
		Eq("status", "pending").
		Execute()

	if err != nil {
		log.Printf("Failed to accept proposal: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to update proposal",
		})
	}

	_, _, err = db.SB.From("tasks").
		Update(map[string]interface{}{"status": "closed"}, "", "").
		Eq("id", proposal.TaskID).
		Execute()

	if err != nil {
		log.Printf("Failed to close task %s: %v", proposal.TaskID, err)
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
		Eq("task_id", proposal.TaskID).
		ExecuteTo(&chats)

	if err != nil {
		log.Printf("Chat lookup DB error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to find chat for task",
		})
	}

	if len(chats) == 0 {
		log.Printf("CRITICAL: No chat found for task %s after proposal %s accepted",
			proposal.TaskID, proposalID)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "No chat found for this task",
		})
	}

	// Добавляем ОБОИХ пользователей в чат: заказчика и исполнителя

	// 1. Добавляем исполнителя (того, кто откликнулся)
	_, _, err = db.SB.From("chat_users").
		Insert(map[string]interface{}{
			"chat_id": chats[0].ID,
			"user_id": proposal.UserID,
		}, false, "", "", "").
		Execute()

	if err != nil {
		log.Printf("CRITICAL: Failed to add freelancer %s to chat %s: %v",
			proposal.UserID, chats[0].ID, err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to add freelancer to chat",
		})
	}

	// 2. Добавляем заказчика (владельца задачи)
	_, _, err = db.SB.From("chat_users").
		Insert(map[string]interface{}{
			"chat_id": chats[0].ID,
			"user_id": tasks[0].ClientID,
		}, false, "", "", "").
		Execute()

	if err != nil {
		log.Printf("CRITICAL: Failed to add client %s to chat %s: %v",
			tasks[0].ClientID, chats[0].ID, err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to add client to chat",
		})
	}

	log.Printf("Added both users to chat %s: freelancer=%s, client=%s",
		chats[0].ID, proposal.UserID, tasks[0].ClientID)

	// --- 11. Отклоняем остальные proposals по этой задаче ---
	_, _, err = db.SB.From("propals").
		Update(map[string]interface{}{"status": "denied"}, "", "").
		Eq("task_id", proposal.TaskID).
		Neq("id", proposalID).
		Execute()

	if err != nil {
		// Не фатально, но логируем как критическое
		log.Printf("CRITICAL: Failed to deny other proposals for task %s: %v",
			proposal.TaskID, err)
		// Не возвращаем 500 — основная операция уже прошла успешно
	}

	log.Printf("Proposal %s accepted successfully for task %s", proposalID, proposal.TaskID)
	return c.JSON(fiber.Map{
		"message": "Proposal accepted and task status updated",
	})
}
