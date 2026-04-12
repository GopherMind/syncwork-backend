package main

import (
	"log"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/Router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db.InitSupabase()
	server := fiber.New()
	Router.AuthRoutes(server)

	server.Use(logger.New())
	server.Use(recover.New())

	log.Fatal(server.Listen(":3000"))
}
