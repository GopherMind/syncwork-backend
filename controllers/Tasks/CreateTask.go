package tasks

import (
	"strings"

	"github.com/GopherMind/syncwork-backend/db"
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

	idUserRaw := c.Locals("id_user")
	idUser, ok := idUserRaw.(int)
	if !ok || idUser == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	createData := map[string]interface{}{
		"title":       task.Title,
		"description": task.Description,
		"budget":      task.Budget,
		"stack":       task.Stack,
		"client_id":   idUser,
	}
	_, _, err = db.SB.From("tasks").Insert(createData, false, "", "", "").Execute()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create task"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Task created successfully", "task": task})
}

// пример запроса: http://localhost:3000/tasks/createTask
// тело запроса:
// {
// 	"title": "Нужен разработчик для создания сайта",
// 	"description": "Ищу опытного разработчика для создания сайта на React. Требуется адаптивный дизайн и интеграция с API.",
// 	"budget": 500,
// 	"stack": ["React", "Node.js", "CSS"]
// }
//
