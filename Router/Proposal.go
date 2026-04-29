package Router

import (
	"github.com/GopherMind/syncwork-backend/controllers/proposals"
	"github.com/GopherMind/syncwork-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProposalRoute(app *fiber.App) {
	group := app.Group("/proposal")

	group.Post("/create/:task", middleware.AuthMiddleware, middleware.RateLimitMiddleware(30), propasals.CreateProposal)
	group.Post("/deny/:id", middleware.AuthMiddleware, propasals.DeniedProposal)
	group.Post("/accept/:id", middleware.AuthMiddleware, propasals.AcceptProposal)
	group.Get("/task", middleware.AuthMiddleware, propasals.GetTaskProposals)
}
