package types

import "github.com/gofiber/fiber/v2"

type BasicResponse struct {
	Success bool        `json:"success"`
	Message *string     `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(BasicResponse{Success: true, Data: data})
}

func Fail(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(BasicResponse{Success: false, Message: &msg})
}
