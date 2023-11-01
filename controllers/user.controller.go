package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wpcodevo/golang-fiber-jwt-rs256/models"
)

// GetMe retrieves the user's record from the context object and returns it to the client in JSON format.
func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}
