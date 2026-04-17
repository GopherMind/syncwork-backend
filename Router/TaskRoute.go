package Router

import (
	tasks "github.com/GopherMind/syncwork-backend/controllers/Tasks"
	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(app *fiber.App) {
	group := app.Group("/tasks")

	group.Get("/getTasks", tasks.GetTasks)
}
// пример запроса: http://localhost:3000/tasks/getTasks?limit=10&price_min=100&price_max=500