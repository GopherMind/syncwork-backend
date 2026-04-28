package Router

import (
	 "github.com/GopherMind/syncwork-backend/controllers/Auth"
	"github.com/GopherMind/syncwork-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	group := app.Group("/auth")

	group.Post("/login", middleware.RateLimitMiddleware(), auth.Register)
	group.Post("/signin", middleware.RateLimitMiddleware(), auth.Signin)
	group.Get("/profile", middleware.AuthMiddleware, auth.GetProfile)
}
