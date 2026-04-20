package web

import (
	"os"
	"path/filepath"
	"strings"

	"cth.release/common"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cfg *common.Config) *fiber.App {
	api := app.Group("/api", EmptyMiddleware)

	api.Get("/health", HealthHandler)

	if hasFile(cfg.Web.StaticDir, cfg.Web.SPAIndexFile) {
		app.Static("/", cfg.Web.StaticDir)
		app.Get("*", SPARouteHandler(cfg))
	}

	return app
}

func EmptyMiddleware(c *fiber.Ctx) error {
	return c.Next()
}

func HealthHandler(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}

func SPARouteHandler(cfg *common.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}
		if strings.HasPrefix(c.Path(), "/api") {
			return c.Next()
		}
		if filepath.Ext(c.Path()) != "" {
			return c.Next()
		}

		return c.SendFile(filepath.Join(cfg.Web.StaticDir, cfg.Web.SPAIndexFile))
	}
}

func hasFile(dir, file string) bool {
	if strings.TrimSpace(dir) == "" || strings.TrimSpace(file) == "" {
		return false
	}

	info, err := os.Stat(filepath.Join(dir, file))
	return err == nil && !info.IsDir()
}
