package main

import (
	"log"

	"cth.release/common"
	"cth.release/web"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := common.GetConfig()
	app := fiber.New()
	web.SetupRoutes(app, config)

	err := app.Listen(":" + config.Port)
	if err != nil {
		log.Fatalf("Error Starting Server: %v", err)
	}
}
