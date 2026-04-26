package tasks

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/models"
	"github.com/GopherMind/syncwork-backend/utils/moderate"
	"github.com/gofiber/fiber/v2"
)

func CreateTask(c *fiber.Ctx) error {
	var task models.FullTask
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if task.Title == "" || task.Description == "" || task.Budget <= 0 || len(task.Stack) == 0 || len(task.Stack) > 10 || len(task.Title) > 100 || len(task.Description) > 2000 || task.Budget > 100000 || task.Level == "" || (task.Level != "junior" && task.Level != "middle" && task.Level != "senior") || task.WorkTimeInWeek < 1 || task.WorkTimeInWeek > 120 {
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

	idUserRaw := c.Locals("user_id")
	idUser, ok := idUserRaw.(string)
	if !ok || idUser == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	createData := map[string]interface{}{
		"title":          task.Title,
		"description":    task.Description,
		"budget":         task.Budget,
		"stack":          task.Stack,
		"client_id":      idUser,
		"level":          task.Level,
		"workTimeInWeek": task.WorkTimeInWeek,
		"status":         "open",
	}
	_, _, err = db.SB.From("tasks").Insert(createData, false, "", "", "").Execute()
	if err != nil {
		fmt.Println("Error inserting task:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create task"})
	}

	body, _, err := db.SB.From("tasks").
		Select("id", "", false).
		Eq("client_id", idUser).
		Limit(1, "").
		Execute()
	if err != nil {
		fmt.Println("Error fetching task:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch task"})
	}

	var insertedTasks []models.FullTask
	if err := json.Unmarshal(body, &insertedTasks); err != nil || len(insertedTasks) == 0 {
		fmt.Println("Error unmarshaling task:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch task"})
	}
	taskID := insertedTasks[0].ID

	DataChat := map[string]interface{}{
		"task_id": taskID,
		"status":  "active",
	}
	_, _, err = db.SB.From("chats").Insert(DataChat, false, "", "", "").Execute()
	if err != nil {
		fmt.Printf("Failed to create chat for task %s: %v\n", task.ID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create chat"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Task created successfully", "task": task})
}
