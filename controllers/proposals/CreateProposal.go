package propasals

import (
	"fmt"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/models"
	"github.com/gofiber/fiber/v2"
)

func CreateProposal(c *fiber.Ctx) error {
	idUserRaw := c.Locals("user_id")
	taskId := c.Params("task")
	idUser, ok := idUserRaw.(string)
	if !ok || idUser == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var proposalBody models.Proposal

	if err := c.BodyParser(&proposalBody); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "fetching data error"})
	}

	if proposalBody.CoverLetter == "" {
		return c.Status(401).JSON(fiber.Map{"error": "cover letter is empty"})
	}

	proposalBody.UserID = idUser
	proposalBody.TaskID = taskId
	proposalBody.Status = "pending"
	
	 var tasks []models.FullTask
	_, err := db.SB.From("tasks").Select("id, client_id", "", false).Eq("id", taskId).ExecuteTo(&tasks)
	if err != nil || len(tasks) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	task := tasks[0]

	clientID := task.ClientID
	if clientID == idUser {
		return c.Status(403).JSON(fiber.Map{"error": "You cannot apply to your own task"})
	}

	insertData := map[string]interface{}{
		"user_id":      proposalBody.UserID,
		"task_id":      proposalBody.TaskID,
		"cover_letter": proposalBody.CoverLetter,
		"status":       proposalBody.Status,
	}
	_, _, err = db.SB.From("propals").Insert(insertData, false, "", "", "").Execute()

	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create propals"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "succes create proposal"})

}
