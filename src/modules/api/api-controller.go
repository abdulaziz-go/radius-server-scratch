package apiModule

import (
	"radius-server/src/database"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	resp := map[string]string{
		"status": "OK",
	}
	dbHealthy := database.HealthCheck()
	if !dbHealthy {
		resp["status"] = "Database is down"
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
