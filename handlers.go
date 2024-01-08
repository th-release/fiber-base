package main

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) *fiber.App {
	app.Get("/hello", helloWorldHandler)
	return app
}

func helloWorldHandler(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}
