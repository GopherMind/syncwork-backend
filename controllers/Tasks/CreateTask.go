package tasks

import (
	"strings"

	"github.com/GopherMind/syncwork-backend/models"
	"github.com/GopherMind/syncwork-backend/utils/moderate"
	"github.com/gofiber/fiber/v2"
)

func CreateTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if task.Title == "" || task.Description == "" || task.Budget <= 0 || len(task.Stack) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Missing required fields: title, description, budget, stack",
			"details": "Ensure title, description, budget , and stack are provided",
		})
	}
	stackStr := strings.Join(task.Stack, ", ")

	contentToModerate := strings.Join([]string{
		task.Title,
		task.Description,
		stackStr,
	}, " | ")
	isSafe, err := moderate.ModerateWithAI(contentToModerate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to moderate task content"})
	}
	if !isSafe {
		return c.Status(400).JSON(fiber.Map{"error": "Task content is not allowed"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Task created successfully", "task": task})
}
