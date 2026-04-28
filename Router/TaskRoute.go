package Router

import (
    "github.com/GopherMind/syncwork-backend/controllers/Tasks"
	"github.com/GopherMind/syncwork-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(app *fiber.App) {
	group := app.Group("/tasks")

	group.Get("/getTasks", tasks.GetTasks)
	group.Post("/createTask", middleware.AuthMiddleware, middleware.RateLimitMiddleware(), tasks.CreateTask)
	group.Get("/getTask/:id", tasks.GetSingleTask)
}

