package propasals

import (
	"fmt"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/gofiber/fiber/v2"
)

func DeniedProposal(c *fiber.Ctx) error {
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
		ID    string `json:"id"`
		Tasks struct {
			ClientID string `json:"client_id"`
		} `json:"tasks"`
	}

	_, err := db.SB.From("propals").
		Select("id, tasks(client_id)", "", false).
		Eq("id", proposalID).
		Eq("status", "pending").
		ExecuteTo(&result)
	if err != nil || len(result) == 0 {
		fmt.Println(err)
		return c.Status(404).JSON(fiber.Map{
			"error":   "NOT_FOUND",
			"message": "Proposal not found",
		})
	}

	if result[0].Tasks.ClientID != userID.(string) {
		return c.Status(403).JSON(fiber.Map{
			"error":   "FORBIDDEN",
			"message": "Only task owner can deny proposals",
		})
	}

	_, _, err = db.SB.From("propals").
		Update(map[string]interface{}{"status": "denied"}, "", "").
		Eq("id", proposalID).
		Eq("status", "pending").
		Execute()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to update proposal",
		})
	}

	return c.JSON(fiber.Map{"message": "Proposal denied"})
}
