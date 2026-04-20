package types

import "github.com/gofiber/fiber/v2"

type BasicResponse struct {
	Success bool        `json:"success"`
	Message *string     `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(BasicResponse{Success: true, Data: data})
}

// Fail — 실패 응답
func Fail(c *fiber.Ctx, status int, msg *string) error {
	return c.Status(status).JSON(BasicResponse{Success: false, Message: msg})
}
