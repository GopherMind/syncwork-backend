package Router

import (
	"github.com/GopherMind/syncwork-backend/controllers/Auth"
	"github.com/gofiber/fiber/v2"
	
)

func AuthRoutes(app *fiber.App) {
	group := app.Group("/auth")

	group.Post("/login", auth.Register)
}