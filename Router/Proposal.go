package Router

import (
	"github.com/GopherMind/syncwork-backend/controllers/proposals"
	"github.com/GopherMind/syncwork-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProposalRoute(app *fiber.App) {
	group := app.Group("/proposal")

	group.Post("/create/:task", middleware.AuthMiddleware, propasals.CreateProposal)
}
