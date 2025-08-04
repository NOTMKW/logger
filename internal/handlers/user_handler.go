package handlers

import (
	"time"

	"github.com/notmkw/logger/internal/services"
	"github.com/notmkw/logger/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	requestID := c.Locals("request_id").(string)

	user, err := h.userService.GetUser(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error(), requestID)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"user_id":    user.ID,
		"request_id": requestID,
		"message":    "User retrieved successfully",
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	requestID := c.Locals("request_id").(string)

	user, err := h.userService.CreateUser()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), requestID)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":    "User created successfully",
		"user_id":    user.ID,
		"request_id": requestID,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	requestID := c.Locals("request_id").(string)

	user, err := h.userService.UpdateUser(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error(), requestID)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"user_id":    user.ID,
		"request_id": requestID,
		"message":    "User updated successfully",
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	requestID := c.Locals("request_id").(string)

	err := h.userService.DeleteUser(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error(), requestID)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":    "User deleted successfully",
		"user_id":    userID,
		"request_id": requestID,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *UserHandler) HealthCheck(c *fiber.Ctx) error {
	requestID := c.Locals("request_id").(string)

	return utils.SuccessResponse(c, fiber.Map{
		"message":    "API is running",
		"request_id": requestID,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
}
