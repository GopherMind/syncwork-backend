package main

import (
	"log"

	"github.com/GopherMind/syncwork-backend/db"
	"github.com/GopherMind/syncwork-backend/Router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	godotenv.Load()

	db.InitSupabase()
	server := fiber.New()
	server.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:5173",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
        AllowCredentials: true,
    }))
	server.Use(logger.New())
	server.Use(recover.New())
	Router.AuthRoutes(server)
	

	log.Fatal(server.Listen(":3000"))
}
