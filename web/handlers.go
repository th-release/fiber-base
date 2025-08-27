package web

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) *fiber.App {
	api := app.Group("/api", EmptyMiddleware)

	api.Get("/health", HealthHandler)
	return app
}

func EmptyMiddleware(c *fiber.Ctx) error {
	return c.Next()
}

func HealthHandler(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}
