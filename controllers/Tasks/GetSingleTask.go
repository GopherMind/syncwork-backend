package tasks

import (
	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/models"
	"github.com/gofiber/fiber/v2"
)
func GetSingleTask(c *fiber.Ctx) error {
	taskID := c.Params("id")

	var tasks []models.FullTask
	if _, err := db.SB.From("tasks").Select("*, profiles(name)", "", false).Eq("id", taskID).ExecuteTo(&tasks); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch task", "details": err.Error()})
	}

	if len(tasks) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	task := tasks[0]

	var propals []map[string]interface{}
	if _, err := db.SB.From("propals").Select("*", "", false).Eq("task_id", taskID).ExecuteTo(&propals); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch proposals", "details": err.Error()})
	}
	task.Proposals = len(propals)

	return c.Status(200).JSON(task)
}

