package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main(){
	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())

	log.Fatal(app.Listen(":3000"))
}