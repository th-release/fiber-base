package main

import (
	"log"

	"cth.release/common"
	"cth.release/common/utils/xss"
	"cth.release/web"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := common.GetConfig()
	xss.ConfigureFromAppConfig(config)
	app := fiber.New()
	web.SetupRoutes(app, config)

	err := app.Listen(":" + config.Server.Port)
	if err != nil {
		log.Fatalf("Error Starting Server: %v", err)
	}
}
