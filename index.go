package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main(){
	//server
	server := fiber.New()

	server.Use(logger.New())
	server.Use(recover.New())

	log.Fatal(server.Listen(":3000"))
}