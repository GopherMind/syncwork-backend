package Router

import (
	tasks "github.com/GopherMind/syncwork-backend/controllers/Tasks"
	"github.com/GopherMind/syncwork-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(app *fiber.App) {
	group := app.Group("/tasks")

	group.Get("/getTasks", tasks.GetTasks)
	group.Post("/createTask", middleware.AuthMiddleware, tasks.CreateTask)
	group.Get("/getTask/:id", tasks.GetSingleTask)
}

// пример запроса: http://localhost:3000/tasks/getTasks?limit=10&price_min=100&price_max=500
// 