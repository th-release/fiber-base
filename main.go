package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	SetupRoutes(app)

	err := app.Listen(":3000")
	if err != nil {
		log.Fatalf("Error Starting Server: %v", err)
	}
}
