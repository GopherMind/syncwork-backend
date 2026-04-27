package propasals

import (
	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/models"
	"github.com/gofiber/fiber/v2"
)

type GroupedProposals struct {
	TaskID    string            `json:"task_id"`
	Proposals []models.Proposal `json:"proposals"`
}

func GetTaskProposals(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(401).JSON(fiber.Map{
			"error":   "UNAUTHORIZED",
			"message": "User ID not found in context",
		})
	}

	var proposals []models.Proposal

	_, err := db.SB.From("propals").
		Select("id, task_id, user_id, cover_letter, status, tasks!inner(client_id)", "", false).
		Eq("tasks.client_id", userID).
		ExecuteTo(&proposals)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_SERVER_ERROR",
			"message": "Failed to fetch proposals: " + err.Error(),
		})
	}

	tempMap := make(map[string][]models.Proposal)
	for _, p := range proposals {
		tempMap[p.TaskID] = append(tempMap[p.TaskID], p)
	}

	result := make([]GroupedProposals, 0, len(tempMap))
	for taskID, props := range tempMap {
		result = append(result, GroupedProposals{
			TaskID:    taskID,
			Proposals: props,
		})
	}

	return c.JSON(result)
}

