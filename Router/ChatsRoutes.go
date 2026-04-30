package Router

import (
	"github.com/GopherMind/syncwork-backend/controllers/chats"
	"github.com/GopherMind/syncwork-backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func ChatsRoutes(app *fiber.App) {
	group := app.Group("/chats")

	group.Get("/getChats", middleware.AuthMiddleware, chats.GetChats)
	group.Get("/messages/:id", middleware.AuthMiddleware, chats.GetChatMessages)
	group.Post("/createMessage/:id", middleware.AuthMiddleware, middleware.RateLimitMiddleware(5), chats.CreateMessage)

	// WebSocket для реалтайм сообщений
	group.Get("/ws/:id", websocket.New(chats.GetChatMessagesWS))
}
