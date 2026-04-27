package tasks

import (
	"strconv"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetTasks(c *fiber.Ctx) error {
	var query models.TaskQuery
	var tasks []models.Task

	if err := c.QueryParser(&query); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid query format"})
	}

	if query.Limit == 0 {
		query.Limit = 6
	}

	q := db.SB.From("tasks").Select("*, profiles(name)", "", false).Limit(query.Limit, "").Eq("status", "open")

	if query.PriceMin > 0 {
		q = q.Gte("budget", strconv.Itoa(query.PriceMin))
	}
	if query.PriceMax > 0 {
		q = q.Lte("budget", strconv.Itoa(query.PriceMax))
	}
	if query.Search != "" {
		q = q.Or("title.ilike.%"+query.Search+"%,description.ilike.%"+query.Search+"%", "")
	}
	if _, err := q.ExecuteTo(&tasks); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tasks", "details": err.Error()})
	}

	return c.JSON(tasks)
}
