package Router

import (
	"github.com/GopherMind/syncwork-backend/controllers/chats"
	"github.com/GopherMind/syncwork-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func ChatsRoutes(app *fiber.App) {
	group := app.Group("/chats")

	group.Get("/getChats", middleware.AuthMiddleware, chats.GetChats)
}
