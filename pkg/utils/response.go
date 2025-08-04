package utils

import "github.com/gofiber/fiber/v2"

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(data)
}

func ErrorResponse(c *fiber.Ctx, status int, message string, requestID string) error {
	return c.Status(status).JSON(fiber.Map{
		"error":      true,
		"message":    message,
		"request_id": requestID,
	})
}

func CreatedResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

func NoContentResponse(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}