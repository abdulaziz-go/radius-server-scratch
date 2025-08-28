package middleware

import (
	"radius-server/src/config"

	"github.com/gofiber/fiber/v2"
)

func CheckXApiKey(c *fiber.Ctx) error {
	xApiKey := c.Get("X-API-KEY")
	if xApiKey == "" || xApiKey != config.AppConfig.Security.XApiKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
			"data":    nil,
		})
	}
	return c.Next()
}
